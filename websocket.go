package poloniex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Channel IDs
const (
	ACCOUNT = 1000 // Account Notification
	TICKER  = 1002 // Ticker
)

const SUBSBUFFER = 256 // Subscriptions Buffer

var (
	channelsByName = make(map[string]int) // channels map by name
	channelsByID   = make(map[int]string) // channels map by id
	marketChannels []int                  // channels list
)

// WSClient describe single websocket connection.
type WSClient struct {
	key        string
	secret     string
	observer   OrderObserver
	Subs       map[string]chan interface{} // subscriptions map
	wsConn     *websocket.Conn             // websocket connection
	wsMutex    *sync.Mutex                 // prevent race condition for websocket RW
	sync.Mutex                             // embedded mutex
}

// NewPublicWSClient creates new web socket public client.
func NewPublicWSClient() *WSClient {
	return &WSClient{
		Subs:    make(map[string]chan interface{}),
		wsMutex: &sync.Mutex{},
	}
}

// NewPrivateWSClient creates new web socket private client.
func NewPrivateWSClient(observer OrderObserver, key, secret string) *WSClient {
	return &WSClient{
		key:      key,
		secret:   secret,
		observer: observer,
		Subs:     make(map[string]chan interface{}),
		wsMutex:  &sync.Mutex{},
	}
}

// Run is connection client to poloniex websocket and start handling messages.
func (ws *WSClient) Run() (err error) {
	logger.Info("connecting to poloniex websocket")

	dialer := &websocket.Dialer{
		HandshakeTimeout: time.Minute,
	}

	wsConn, _, err := dialer.Dial(pushAPIUrl, nil)
	if err != nil {
		return
	}

	ws.wsConn = wsConn

	if err = setChannelsID(); err != nil {
		return
	}

	go func() {
		for {
			err := ws.wsHandler()
			if err != nil {
				wsConn, _, _ := dialer.Dial(pushAPIUrl, nil)
				ws.wsConn = wsConn
			}
		}
	}()

	logger.Info("successfully connected to poloniex")

	return nil
}

// Web socket reader.
func (ws *WSClient) readMessage() ([]byte, error) {
	ws.wsMutex.Lock()
	defer ws.wsMutex.Unlock()
	_, rmsg, err := ws.wsConn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return rmsg, nil
}

// Web socket writer.
func (ws *WSClient) writeMessage(msg []byte) error {
	ws.wsMutex.Lock()
	defer ws.wsMutex.Unlock()
	return ws.wsConn.WriteMessage(1, msg)
}

func setChannelsID() (err error) {
	publicAPI := NewPublicClient()

	tickers, err := publicAPI.GetTickers()
	if err != nil {
		return err
	}

	for k := range tickers {
		id := tickers[k].ID
		channelsByName[k] = id
		channelsByID[id] = k
		marketChannels = append(marketChannels, id)
	}

	channelsByName["TICKER"] = TICKER
	channelsByID[TICKER] = "TICKER"

	channelsByName["ACCOUNT"] = ACCOUNT
	channelsByID[ACCOUNT] = "ACCOUNT"

	return
}

// Create handler.
// If the message comes from the channels that are subscribed,
// it is sent to the chans.
func (ws *WSClient) wsHandler() error {
	for {
		msg, err := ws.readMessage()
		if err != nil {
			return err
		}

		var imsg []interface{}
		err = json.Unmarshal(msg, &imsg)
		if err != nil || len(imsg) < 3 {
			continue
		}

		arg, ok := imsg[0].(float64)
		if !ok {
			continue
		}

		chID := int(arg)
		args, ok := imsg[2].([]interface{})
		if !ok {
			continue
		}

		var wsUpdate interface{}

		switch {
		case chID == TICKER:
			wsUpdate, err = convertArgsToTicker(args)
			if err != nil {
				logger.WithError(err).Error("can not parse ticker message")
				continue
			}
		case chID == ACCOUNT:
			wsUpdate, err = convertArgsToAccountNotification(args)
			if err != nil {
				logger.WithError(err).Error("can not parse account notification message")
				continue
			}
		case intInSlice(chID, marketChannels):
			wsUpdate, err = convertArgsToMarketUpdate(args)
			if err != nil {
				logger.WithError(err).Error("can not parse market update message")
				continue
			}
		default:
			continue
		}

		chName := channelsByID[chID]
		if ws.Subs[chName] != nil {
			select {
			case ws.Subs[chName] <- wsUpdate:
			default:
			}
		}
	}
}

// sub-function for subscription.
func (ws *WSClient) subscribe(chID int, chName string) (err error) {
	ws.Lock()
	defer ws.Unlock()

	if ws.Subs[chName] == nil {
		ws.Subs[chName] = make(chan interface{}, SUBSBUFFER)
	}

	subsMsg, ok := subscription{
		Command: "subscribe",
		Channel: strconv.Itoa(chID),
	}.toJSON()
	if !ok {
		return errors.New("failed to convert subscription struct to JSON")
	}

	err = ws.writeMessage(subsMsg)
	if err != nil {
		return
	}

	return
}

// sub-function for unsubscription.
// the chans are not closed once the subscription is made to protect chan address.
// To prevent chans taking a new address on the memory, thus chans can be used repeatedly.
func (ws *WSClient) unsubscribe(chName string) (err error) {
	ws.Lock()
	defer ws.Unlock()

	if ws.Subs[chName] == nil {
		return
	}

	unSubsMsg, _ := subscription{
		Command: "unsubscribe",
		Channel: chName,
	}.toJSON()

	err = ws.writeMessage(unSubsMsg)
	if err != nil {
		return err
	}

	return
}

func (ws *WSClient) sign(formData string) (signature string, err error) {
	if ws.key == "" || ws.secret == "" {
		err = Error(SetAPIError)
		return
	}

	mac := hmac.New(sha512.New, []byte(ws.secret))
	_, err = mac.Write([]byte(formData))
	if err != nil {
		return
	}

	signature = hex.EncodeToString(mac.Sum(nil))
	return
}
