//TODO: fix const time as it is not formatted correctly and add where time/date is needed
const time = () => { return new Date().toLocaleString() };
let uName = ""
class MySocket {
  wsType = ""
  constructor() {
    this.mysocket = null;
  }
  // TODO: insert user ID variable, participants needs to be filled
  sendNewChatRequest() {
    console.log("new chat request")
    let m = {
      type: 'chat',
      timestamp: time(),
      conversations: [
        {
          participants: [
            //sender: bar userID
            {
              id: "975496ca-9bfc-4d71-8736-da4b6383a575",
            },
            //other participants (receiver): foo userID
            {
              id: "6d01e668-2642-4e55-af73-46f057b731f9",
            }
          ],
          chats: [
            {
              sender: {
                // TODO: this is just the first placeholder above, once the user is logged in and their ID is stored client side this ID should represent the logged in user
                // bar userID
                id: "975496ca-9bfc-4d71-8736-da4b6383a575",
              },
              body: document.getElementById('chatIPT').value,
            }
          ]
        }
      ]
    }
    this.mysocket.send(JSON.stringify(m));
    document.getElementById('chatIPT').value = ""
  }

  sendChatContentRequest(e, chat_id = "") {
    this.mysocket.send(JSON.stringify({
      type: "content",
      resource: e.target.id,
      chat_id: chat_id,
    }));
  }
  getClickedParticipantID() {
  }
  getLoggedInUserID() {
  }

  
  keypress(e) {
    if (e.keyCode == 13) {
      this.wsType = e.target.id.slice(0, -3)
      if (this.wsType = 'chat') {
        this.sendNewChatRequest()
      }
    }
  }
  // registerHandler(text){
  //   console.log("register handler")
  // }
  chatHandler(text) {
    const m = JSON.parse(text)
    for (let c of m.conversations) {
      for (let p of c.chats) {
        let chat = document.createElement("div");
        chat.className = "submittedchat"
        chat.id = p.chat_id
        chat.innerHTML = "<b>Me: " + p.sender.nickname + "</b>" + "<br>" + "<b>Date: " + "</b>" + p.date + "<br>" + p.body + "<br>";
        document.getElementById("chatcontainer").appendChild(chat)
      }
    }
  }
  contentHandler(text) {
    const c = JSON.parse(text)
    document.getElementById("content").innerHTML = c.body;
  }
  presenceHandler(text) {
    console.log(text)
    const m = JSON.parse(text)

    let presenceCont = document.getElementById("presencecontainer")
    if (presenceCont.childElementCount != 0) {
      while (presenceCont.firstChild) {
      presenceCont.removeChild(presenceCont.lastChild);
    }
  }
    if (m.presences != null ) {
      for (let p of m.presences) {
      const consp = p
      let user = document.createElement("button");
      user.addEventListener('click', function (event, chat = consp) {
        event.target.id = "chat"
        contentSocket.sendChatContentRequest(event, chat.chat_id)
      });
      user.id = p.id
      user.innerHTML = p.nickname
      user.style.color = 'white'
      user.className = "presence " + p.nickname 
      presenceCont.appendChild(user)
    }
    }
    console.log("Presences successfully updated")
  }
  postHandler(text) {
    console.log("Printing posts....")
    const m = JSON.parse(text)
    console.log(document.getElementById("submittedposts"))
    for (let p of m.posts) {
      console.log(document.getElementById("submittedposts"))
      const consp = p
      let post = document.createElement("div");
      post.className = "submittedpost"
      post.id = p.post_id
      post.innerHTML = "<b>Title: " + p.title + "</b>" + "<br>" + "<b>Nickname: " + "</b>" + p.nickname + "<br>" + "<b>Category: " + p.categories + "</b>" + "<br>" + p.body + "<br>";
      let button = document.createElement("button")
      button.classname = "addcomment"
      button.innerHTML = "Comments"
      button.addEventListener('click', function (event, post = consp) {
        event.target.id = "comment"
        contentSocket.sendContentRequest(event, post.post_id)
      });
      post.appendChild(button)

      console.log(document.getElementById("submittedposts"))
      if ( document.getElementById("submittedposts") != null){
        console.log("appending Postss")
      document.getElementById("submittedposts").appendChild(post)
      } else {
        console.log("Submitted posts == null")
      }
    }
  }
  sendNewCommentRequest(e) {
    console.log("adding comments",uName)
    let post = document.getElementById('postcontainerforcomments')
    for (const child of post.children) {
      if (containsNumber(child.id)) {
        let m = {
          type: 'post',
          timestamp: time(),
          posts: [
            {
              post_id: child.id,
              comments: [
                {
                  post_id: child.id,
                  nickname: uName,
                  body: document.getElementById('commentbody').value,
                }
              ]
            }
          ]
        }
        this.mysocket.send(JSON.stringify(m));
        document.getElementById('commentbody').value = ""
      }
    }
  }
  sendNewPostRequest(e) {
    console.log("sending new post request",uName)
    let m = {
      type: 'post',
      timestamp: time(),
      posts: [
        {
          //nickname: e.target.nickname,
          nickname: uName,
          title: document.getElementById('posttitle').value,
          categories: document.getElementById('category').value,
          body: document.getElementById('postbody').value,
        }
      ]
    }
    this.mysocket.send(JSON.stringify(m));
    document.getElementById('posttitle').value = ""
    document.getElementById('category').value = ""
    document.getElementById('postbody').value = ""
  }
  sendSubmittedPostsRequest() {
    console.log("Sending post request...")
    this.mysocket.send(JSON.stringify({
      type: "post",
    }));
  }
  sendContentRequest(e, post_id = "") {
    console.log("Sending content request...")
    this.mysocket.send(JSON.stringify({
      type: "content",
      resource: e.target.id,
      post_id: post_id,
    }));
  }
  sendChatContentRequest(e, chat_id = "") {
    this.mysocket.send(JSON.stringify({
      type: "content",
      resource: e.target.id,
      chat_id: chat_id,
    }));
  }
 
sendPresenceRequest() {
  console.log("Updating Presences....")
  this.mysocket.send(JSON.stringify({
    type: "presence",
    username: uName,
  }));
}
  connectSocket(URI, handler) {
    if (URI === 'chat') {
      this.wsType = 'chat'
      console.log("Chat Websocket Connected");
    }
    if (URI === 'content') {
      this.wsType = 'content'
      console.log("Content Websocket Connected");
    }
    if (URI === 'post') {
      this.wsType = 'post'
      console.log("Post Websocket Connected");
    }
    if (URI === 'presence') {
      this.wsType = 'presence'
      console.log("Presence Websocket Connected");
    }
    var socket = new WebSocket("ws://localhost:8080/" + URI);
    this.mysocket = socket;
    socket.onmessage = (e) => {
      // console.log("socket message")
      handler(e.data, false);
    };
    socket.onopen = () => {
      // console.log("socket opened");
    };
    socket.onclose = () => {
      // console.log("socket closed");
    };
  }

}
function containsNumber(str) {
  return /[0-9]/.test(str);
}
  
//object to store form data
let registerForm = {
  nickname : "",
  age: "",
  gender: "",
  fName: "",
  lName: "",
  email: "",
  password: "",
  loggedin: "false",
}

let loginForm ={
  nickname:"",
  password:"",
}

//******************* */gets registration form details*******************************
function getRegDetails(){

    //creates array of gender radio buttons 
  let genderRadios = Array.from (document.getElementsByName('gender'))
  for(let i=0; i <genderRadios.length; i ++){
    // console.log(genderRadios[i].checked)
    if(genderRadios[i].checked){ //stores checked value
      registerForm.gender = genderRadios[i].value
    }
  }
// POPULATE REGISTER FORM WITH FORM VALUES
    registerForm.nickname = document.getElementById('nickname').value
    uName = registerForm.nickname 
    registerForm.age = document.getElementById('age').value
    registerForm.firstname = document.getElementById('fname').value
    registerForm.lastname = document.getElementById('lname').value
    registerForm.email = document.getElementById('email').value
    registerForm.password = document.getElementById('password').value
    //convert data to JSON
    let jsonRegForm = JSON.stringify(registerForm)
    // console.log(jsonRegForm)
    if(registerForm.password.length <5){
      registerForm.password = ""
    }
    
// SEND DATA TO BACKEND USING FETCH
  console.log(registerForm)
    if(registerForm.nickname !=""&& registerForm.email !="" &&registerForm.password !="" ){
        
    fetch("/register",{
      headers:{
        'Accept':'application/json',
        'Content-Type': 'application/json'
      },
      method: "POST",
      body:jsonRegForm

    }).then((response)=>{
      response.text().then(function (jsonRegForm){
        let result = JSON.parse(jsonRegForm)
        console.log(result)
      })

    }).catch((error)=>{
      console.log(error)
    })

    document.getElementById('register').reset()
    alert("successfully registered")
  }
}


// **********************************LOGIN*******************************************
function loginFormData(){
  loginForm.nickname = document.getElementById('nickname-login').value
  loginForm.password = document.getElementById('password-login').value
  uName = loginForm.nickname
  // console.log(loginForm)

  let loginFormJSON = JSON.stringify(loginForm)
  // console.log(loginFormJSON)
  let logindata = {nickname:"",
                  password:"",}
  // let id = ""

  fetch("/login",{
      headers:{
        'Accept':'application/json',
        'Content-Type': 'application/json'
      },
      method: "POST",
      body:loginFormJSON

    }).then((response)=>{

      response.text().then(function (loginFormJSON){
        // JSON.parse(loginFormJSON)
        let result = JSON.parse(loginFormJSON)
        console.log("parse",result)

        
        if (result == null){
          alert("incorrect username or password")

        } else{
          logindata.nickname = result[0].nickname
          // logindata.password = result[0].password
          user.innerText = `Hello ${document.cookie.match(logindata.nickname)}`
          alert("you are logged in ")
          document.getElementById("login").style.display = "none"
            document.getElementById("logout").style.display="block"
            document.getElementById("profile").style.display="block"
            document.getElementById("postLogin").style.display="none"
            document.getElementById("postButton").style.display="block"
            document.getElementById("postButton").style.margin="0 auto"
          presenceSocket.sendPresenceRequest()
        }
      })

    }).catch((error)=>{
      console.log(error)

    })
    // console.log("logindata",logindata, "hi")
    // console.log( Object.keys(logindata).length)
    // console.log(JSON.stringify(logindata))

  document.getElementById('login-form').reset()

    let user= document.getElementById('welcome')
  // document.getElementById('login-form').reset()
  // console.log(t)

}
class User {
  constructor(nickname, userID) {
    nickname = "";
    userID =  "";
  }
}
let logoutData = {
  username: ""
}
function Logout() {
  let cookies = document.cookie
  let username = (cookies.split("="))[0]
  logoutData.username= username
  console.log(username)
  let logoutDataJSON = JSON.stringify(logoutData)
  fetch("/logout",{
headers:{
'Accept':'application/json',
'Content-Type': 'application/json'
},
method: "POST",
body: logoutDataJSON
}).then((response)=>{
  document.getElementById("login").style.display = "block"
  document.getElementById("logout").style.display="none"
  document.getElementById("profile").style.display="none"
  document.getElementById("postLogin").style.display="block"
  document.getElementById("postButton").style.display="none"
console.log("Logged out", response)
presenceSocket.sendPresenceRequest()
})
let user= document.getElementById('welcome')
user.innerText = "Welcome"
}