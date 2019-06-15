import React, {useState, useEffect} from 'react';
import styled from 'styled-components';
import './App.css';

let socket = new WebSocket("ws://localhost:8080/ws");

function App() {
  const [inRoom, setInRoom] = useState(false);
  const [message, setMessage] = useState('');
  const [user, setUser] = useState('');
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
    let temp = JSON.parse(event.data)
    setMessages([...messages, {...temp}])
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
    setMessages([...messages, {user, message}])
  }
  const retMessages = messages.map(resp => {
    if (resp.user === user) {
      return (
      <Left>
        <User>{resp.user}</User><p>{resp.message}</p>
      </Left>
      )
    } else if (resp.message.includes(resp.user)) {
      return (
      <Right>
        <Message>{resp.message}</Message>
      </Right>
      )
    } else {
      return (
        <Right>
          <User>{resp.user}</User>
          <Message>{resp.message}</Message>
        </Right>
      )
    }
  })

  const keyPress = (e) => {
    if(e.keyCode === 13){
       console.log('value', e.target.value);
       // put the login here
    }
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
      <Container>
          {retMessages}
      </Container>
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
  width: 50%;
`

const User = styled.p`
  align-self: center;
`

const Message = styled.p`
  align-self: flex-end;
`

const Left = styled.div`
  align-self: flex-start;
  display: flex;
  flex-direction: column;
  background-color: red;
  width: 50%;
`

const Right = styled.div`
  display: flex;
  align-self: flex-end;
  flex-direction: column;
  width: 50%;
  background-color: blue;
`

export default App;
