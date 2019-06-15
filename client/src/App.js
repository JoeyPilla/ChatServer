import React, {useState, useEffect} from 'react';
import logo from './logo.svg';
import './App.css';

let socket = new WebSocket("ws://localhost:8080/ws");

function App() {
  const [messageCount, setMessageCount] = useState(0);
  const [inRoom, setInRoom] = useState(false);
  const [message, setMessage] = useState('');
  const [user, setUser] = useState('');

   useEffect(() => {
    if(inRoom) {
      console.log('joining room');
      socket.send(JSON.stringify(
        {
          user: {
            user: user,
            action: "Enter Room"
          },
          message: {
            message: user+" has joined the chat"
          }
        }));
    }

    return () => {
      if(inRoom) {
        socket.send(JSON.stringify(
          {
            user: {
              user: user,
              action: "Leave Room"
            },
            message: {
              message: user+ " has left the chat"
            }
          }));
      }
    } 
  }, [inRoom, , user]);

  useEffect(() => {
    socket.onmessage = () => (msg => {
      console.log(msg);
    });

    document.title = `${messageCount} new messages have been emitted`;
  }, [messageCount]); //only re-run the effect if new message comes in

  const handleInRoom = () => {
    inRoom
      ? setInRoom(false)
      : setInRoom(true);
  }

  const handleNewMessage = () => {
    socket.send(JSON.stringify(
      {
        user: {
          user: user,
          action: "Message"
        },
        message: {
          message: message
        }
      }));
    setMessageCount(messageCount + 1);
  }
  return (
    <div>
      <h1>
        {inRoom && `Inside Room`}
        {!inRoom && `Outside Room` }
      </h1>
      <h4>
      {inRoom && user}
      </h4>  
      {!inRoom &&
              <input
              type="text"
              value={user}
              onChange={(e) => setUser(e.target.value)}
              placeholder="enter username here"
            />
      }

      {inRoom &&
        <>
                <input
                type="text"
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                placeholder="enter message here"
              />
        <button onClick={() => handleNewMessage()}>
          Send New Message
        </button>
        </>
        }
        <button onClick={() => handleInRoom()}>
          {inRoom && `Leave Room` }
          {!inRoom && `Enter Room` }
      </button>
    </div>
  );
}


export default App;
