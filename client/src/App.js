import React, {useState, useEffect, useRef} from 'react';
import styled from 'styled-components';
import './App.css';

let socket = new WebSocket("ws://localhost:8080/ws");


const Messages = ( {messages} ) => {
  const messagesEndRef = useRef(null);
  const scrollToBottom = () => {
    messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
  };
  useEffect(scrollToBottom, [messages]);
  return (
    <Container>
      {messages}
      <div ref={messagesEndRef} />
    </Container>
  );
};




function App() {
  const [inRoom, setInRoom] = useState(false);
  const [message, setMessage] = useState('');
  const [user, setUser] = useState('');
  const [prevUser, setPrevUser] = useState('')
  const [messages, setMessages] = useState([]);
   useEffect(() => {
    if(inRoom) {
      socket.send(JSON.stringify(
        {
          user: {
            user: user,
            action: "Enter Room"
          },
            message: user + " has joined the chat"
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
            message: user + " has left the chat"
          }));
      }
    } 
   }, [inRoom, , user]);
  
  socket.onmessage = function (event) {
    let {user, message} = JSON.parse(event.data)
    if (message.includes(user + " has") || prevUser === user) {
        setMessages([...messages,
        <Right>
          <Message>{message}</Message>
        </Right>
          ])
    } else {
        setPrevUser(user);
        setMessages([...messages,
        <Right>
          <User>{user}</User>
          <Message>{message}</Message>
        </Right>
        ])
    }
  };

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
          message: message
      }));
    setMessages([...messages,
      <Left>
        <SentMessage>{message}</SentMessage>
      </Left>])
  }

  return (
    <Home>
      <div>
        <h1>
          {inRoom && `Inside Room`}
          {!inRoom && `Outside Room` }
        </h1>
        <h4>
          {inRoom && user}
        </h4>
      </div>
      <Messages messages={messages}/>
      <Footer>
        {!inRoom &&
          <input
            type="text"
            value={user}
            onChange={(e) => setUser(e.target.value)}
            placeholder="enter username here"
            onKeyDown={(e) => {
              if (e.keyCode === 13) {
                handleInRoom()
              }
            }}/>
        }
        {inRoom &&
          <>
            <input
              type="text"
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder="enter message here"
              onKeyDown={(e) => {
                if (e.keyCode === 13) {
                  handleNewMessage()
                }
              }}/>
            <button onClick={() => handleNewMessage()}>
              Send New Message
            </button>
          </>
          }
        <button onClick={() => handleInRoom()}>
          {inRoom && `Leave Room` }
          {!inRoom && `Enter Room` }
        </button>
      </Footer>
    </Home>
  );
}

const Home = styled.div`
  display: grid;
  grid-template-columns: 1fr;
  grid-template-rows: 100px 1fr 100px;
  width: 100vw;
  height: 100vh;
`

const Footer = styled.div`
  display:flex;
  justify-content: center
  position: fixed;
  bottom: 0;
  background-color: white;
  width: 100%;
  height: 75px;
`

const Container = styled.div`
  display: flex;
  flex-direction: column;
  justify-self:center;
  width: 100%;
  margin-bottom: 75px;
  max-width: 450px;
`

const User = styled.p`
  align-self: center;
  display: inline-block;
  word-wrap: break-word;
  word-break: break-word;
`

const Message = styled.p`
display: inline-block;
  align-self: flex-end;
  word-wrap: break-word;
  word-break: break-word;
  padding-right: 30px;
`
const SentMessage = styled.p`
display: inline-block;
  align-self: flex-start;
  word-wrap: break-word;
  word-break: break-word;
  padding-left: 30px;
`

const Left = styled.div`
  align-self: flex-start;
  display: flex;
  flex-direction: column;
  width: 50%;
`

const Right = styled.div`
  display: flex;
  align-self: flex-end;
  flex-direction: column;
  width: 50%;
`

export default App;
