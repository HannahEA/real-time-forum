//TODO: fix const time as it is not formatted correctly and add where time/date is needed
const time = () => { return new Date() };
let uName = ""
let sender_id
let reciev_id 
function isValid(date) {
  return !isNaN(date.getTime())
}
class MySocket {
  wsType = ""
  constructor() {
    this.mysocket = null;
  }
  // TODO: insert user ID variable, participants needs to be filled
  sendNewChatRequest() {
    console.log("new chat request")
    console.log(reciev_id)
    console.log(time())
    let m = {
      type: 'chat',
      timestamp: `${time()}`,
      ConvoID: "1234",
      conversations: [
        {
          participants: [
            //sender: bar userID
            {
              id: sender_id,
            },
            //other participants (receiver): foo userID
            {
              id: reciev_id,
            }
          ],
          chats: [
            {
              sender: {
                // TODO: this is just the first placeholder above, once the user is logged in and their ID is stored client side this ID should represent the logged in user
                // bar userID
                id: sender_id,
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

  getAllChats(reciev_id){

   
    // console.log(reciev_id)
    let getCorrectChats ={
      user1: getCookieName(),
      user2: reciev_id,
    }
    let stringified = JSON.stringify(getCorrectChats)
    fetch("/getChats",{
      headers:{
        'Accept':'application/json',
        'Content-Type': 'application/json'
      },
      method: "POST",
      body: stringified
    })
    .then(response => 
      response.json())
    .then(data => chatSocket.text(data))
    .catch(error => console.log(error))
  }

  text(data){

    let m = {
      type: 'chat',
      timestamp: time(),
      conversations:[
        {
          participants: [
            {
              id: `${getCookieName()}`
              // id: `${data.conversations.chats.id}`
            },
            {
              id: `${reciev_id}`
            }
          ],
          chats: data
        }
      ]
      }
      
    chatSocket.chatHandler(m)
  }


  sendChatContentRequest() {
    this.mysocket.send(JSON.stringify({
      type: "content",
      resource: "chat",
    }));
  }

  
  keypress(e) {
    if (e.keyCode == 13) {
      this.wsType = e.target.id.slice(0, -3)
      if (this.wsType = 'chat') {
        this.sendNewChatRequest()
      }
    }
  }

  chatHandler(text) {
   console.log("printing chats..", text)
  //  console.log("type of chat message", typeof text)
   let chatBox = document.getElementById("newchatscontainer")
    if (typeof text == "string") {
      //new message
      let x = JSON.parse(text)
      console.log(x)
      //move presence to top when you get a chat for reciever
      if (x.sender.id == getCookieName()) {
        //push this chat to the top for the sender
        let pres = document.getElementById(reciev_id)
        let newP = pres
        pres.remove()
        document.getElementById("presencecontainer").prepend(newP)
      }
      // is the chat box open?
      
      console.log("chatBox", chatBox)
      if (chatBox != null) {
         // Yes =
        // console.log("new chat message", text)
        let chat = document.createElement("div");
        chat.className = "submittedchat"
        chat.id = x.chat_id
        let date1 = new Date(x.date)
        chat.innerHTML = "<b>Me: " + x.sender.id + "</b>" + "<br>" + "<b>Date: " + "</b>" + reformatTime(date1) + "<br>" + x.body + "<br>";
        document.getElementById("newchatscontainer").appendChild(chat)
      } else {
        //No = send POST fetch with message containing sender and reciever
        x.notification = true
        x.reciever.id = `${getCookieName()}`
        // console.log("notifictaion message", x)
        var stringified = JSON.stringify(x)
        Notifications(stringified)
      }  
     } else{
      //full chat histry
      
      let position 
      let chats = text.conversations[0].chats
      let chatL = text.conversations[0].chats.length-1
      //Get the height of the chat bpox before loading next 10
//let prevHeight =document.getElementById("newchatscontainer").scrollHeight

      for ( position = chatL ; position>= chatL - 9; position--) {
        // console.log("position", position)
        if (position == -1) {
          break
        }
         let p = chats[position]
          let chat = document.createElement("div");
          chat.className = "submittedchat"
          chat.id = p.chat_id
          let date1 = new Date(p.date)
          console.log("date", date1)
          chat.innerHTML = "<b>Me: " + p.sender.id + "</b>" + "<br>" + "<b>Date: " + "</b>" + reformatTime(date1)+"<br>" + p.body + "<br>";
          document.getElementById("newchatscontainer").prepend(chat)
      }
      chatBox.scrollTop = chatBox.scrollHeight;
      position--
      chatBox.addEventListener('scroll', (event)=>{
       
        if (chatBox.scrollTop == 0) {
            for (let i = 0; i<=10;i++ ) { 
              console.log("position2", position)
              if (position == -1) {
                console.log("position is -1")
                break
              }
              let p = chats[position]
              let chat = document.createElement("div");
              chat.className = "submittedchat"
              chat.id = p.chat_id
              let date1 = new Date(p.date)
              chat.innerHTML = "<b>Me: " + p.sender.id + "</b>" + "<br>" + "<b>Date: " + "</b>" + reformatTime(date1) + "<br>" + p.body + "<br>";
              document.getElementById("newchatscontainer").prepend(chat)
              position--
          
          }
          chatBox.scrollTop = chatBox.scrollHeight;
          
          //get the current height of the chatbox 
          //let currentHeight = document.getElementById("newchatscontainer").scrollHeight
          //Scroll from top by currentheight - previous height
          //document.getElementById("newchatscontainer").scrollTo = currentHeight - prevHeight
      }})
     }
    
  }
  contentHandler(text) {
    const c = JSON.parse(text)
    console.log("content", c)
    document.getElementById("content").innerHTML = c.body;
    if (c.resource == "chat") {
      let titleBox = document.getElementById("chatTitle")
      let title = document.createElement("h2")
      title.innerHTML = reciev_id
      title.style.color = "gray"
     
      titleBox.append(title)
    }
    
  }
  presenceHandler(text) {
   
    const m = JSON.parse(text)
    console.log(m.presences)
    //reorder by online and offline
    let onlineUser = []
    let offlineUser = []
    if (m.presences) {
     m.presences =  m.presences.map(({online}, i)=>{ 
    if (online){
      onlineUser.push(m.presences[i])
     }else{
      offlineUser.push(m.presences[i])
     }})
    }
    onlineUser = onlineUser.concat(offlineUser)
  // reorder by latest chat 
  onlineUser.map(function (element, i) {
   var date1= new Date(onlineUser[i].last_contacted_time)
   onlineUser[i].last_contacted_time = date1
   return onlineUser[i]
    
  })
  
  console.log(onlineUser)
    // var date1 = new Date(onlineUser[i].last_contacted_time)
    // console.log( date1)
   
 onlineUser.sort(function(a,b) {
  console.log(typeof a.last_contacted_time)
  let aDate = new Date (a.last_contacted_time)
  let bDate = new Date (b.last_contacted_time)
  if (!isValid(aDate) && !isValid(bDate)) {
    return 0
  }
  if (!isValid(aDate)){
    return 1
  }
  if (!isValid(bDate)){
    return -1
  }
  return bDate-aDate})
  console.log(onlineUser)
 //  console.log("pres in time order", m.presences )
    //remove old list     
    let presenceCont = document.getElementById("presencecontainer")
    if (presenceCont.childElementCount != 0) {
      while (presenceCont.firstChild) {
      presenceCont.removeChild(presenceCont.lastChild);
      }
    }
    console.log("online users", onlineUser)
    if (m.presences != null) {
      //for each user
      for (let p of onlineUser) {
      //create chat button
        let user = document.createElement("button");
        user.addEventListener('click', function ( ) {
          reciev_id = p.nickname
          sender_id = `${getCookieName()}`
          contentSocket.sendChatContentRequest()
          chatSocket.getAllChats(reciev_id)
          
          // check for notifictaion symbol on button if present remove class before requesting all chats 
          if (document.getElementById(p.nickname).style.backgroundColor === "purple" ){
            user.style.backgroundColor = "#9c9494"
            // if (user.classList.contains('offline')){
            //  user.style.backgroundColor = 'red'
            // } else {
            //  user.style.backgroundColor = 'black'
            // }
            let m = {
              reciever: {
                id: reciev_id,
              },
              sender: {
                id: sender_id,
              },
              notification: false
            }
            let stringified = JSON.stringify(m)
            Notifications(stringified)
            console.log("notification removed")
          }
        });
       
        user.id = p.nickname
        user.innerHTML = p.nickname
        user.style.color = 'white'
        console.log(p.online, p.nickname)
        user.className = "presence " + p.nickname
        console.log("p.notification", p.notification)

        if (p.notification == true) {
          user.style.backgroundColor = "purple"
        }
        
        if (p.online === false) {
        //  user.style.backgroundColor = "red"
          user.classList.add('offline')
          user.innerHTML += `<div class="icon">
            <i class="material-icons" style="color:red"></i>
            </div>`
        } else {
              user.classList.add('online')
              user.innerHTML += `<div class="icon">
                <i class="material-icons" style="color:white"></i>
                </div>`
        }

        presenceCont.appendChild(user)

      } 
    }else{
      console.log('empty presence list sent to phandler')
      presenceSocket.sendPresenceRequest()
    }
    console.log("Presences successfully updated")

    
    
  }
  postHandler(text) {
    console.log("Printing posts....")
    const m = JSON.parse(text)
    console.log("post message", m)
    let content = document.getElementById("content")
    let children = content.childNodes;
    console.log("content children", children)
    for (let c of children) {
      if (c.id != "submittedposts") {
        console.log("child", c)
        // let child = document.getElementById(`${c.id}`)
        c.remove()
      }
    }
    
    
    if (m.posts.length != 0) {
      console.log("submitted posts?", document.getElementById("submittedposts"))
      if (document.getElementById("submittedposts") == null) {
        let subPosts = document.createElement('div')
      subPosts.id = "submittedposts"
      content.appendChild(subPosts)
      }
      
    
    for (let p of m.posts) {
      const consp = p
      let post = document.createElement("div");
      let button = document.createElement("button")
      post.className = "submittedpost"
      post.id = p.post_id
      post.innerHTML = "<b>Title: " + p.title + "</b>" + "<br>" + "<b>Nickname: " + "</b>" + p.nickname + "<br>" + "<b>Category: " + p.categories + "</b>" + "<br>" + p.body + "<br>";
      
      button.className = "addcomment"
      button.innerHTML = "Comments"
      button.addEventListener('click', function (event, post = consp) {
        event.target.id = "comment"
        contentSocket.sendContentRequest(event, post.post_id)
      });
      post.appendChild(button)
     document.getElementById("submittedposts").appendChild(post)
    }
    }
    if ( document.getElementById("submittedposts") != null){
      console.log("appended Posts")
     
    } else {
      console.log("Submitted posts == null")
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
                  nickname: `${getCookieName()}`,
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
    if (document.getElementById('posttitle').value != "" &&
    document.getElementById('category').value != "" &&
    document.getElementById('postbody').value != "") {
      let m = {
        type: 'post',
        timestamp: time(),
        posts: [
          {
            //nickname: e.target.nickname,
            nickname: `${getCookieName()}`,
            title: document.getElementById('posttitle').value,
            categories: document.getElementById('category').value,
            body: document.getElementById('postbody').value,
          }
        ]
      }
      this.mysocket.send(JSON.stringify(m));
    }
    
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
  // sendChatContentRequest(e, chat_id = "") {
  //   this.mysocket.send(JSON.stringify({
  //     type: "content",
  //     resource: e.target.id,
  //     chat_id: chat_id,
  //   }));
  // }
 
sendPresenceRequest() {
  console.log("Updating Presences....")
  this.mysocket.send(JSON.stringify({
    type: "presence",
    username: `${getCookieName()}`,
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
  loggedin: false,
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
    if(registerForm.nickname !=""&& checkEmail(registerForm.email) &&registerForm.password !="" ){
        
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
        if (result.Email === "true" || result.Nickname === "true") {
          alert("Nickname or email already exists")
        } else {
          alert("successfully registered")
        }
      })

    }).catch((error)=>{
      console.log(error)
    })
    document.getElementById('register').reset()
    
  }
}

function checkEmail(email) {
  return email.includes("@")
}

// **********************************LOGIN*******************************************
var presenceSocket = new MySocket
var  chatSocket = new MySocket

function loginFormData(event){
  loginForm.nickname = document.getElementById('nickname-login').value
  loginForm.password = document.getElementById('password-login').value
  
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
          uName = result[0].nickname
          // logindata.password = result[0].password
          let user= document.getElementById('welcome')
          user.innerText = `Hello ${logindata.nickname}`
          alert("you are logged in ")
          document.getElementById("login").style.display = "none"
            document.getElementById("logout").style.display="block"
           
            document.getElementById("postLogin").style.display="none"
            document.getElementById("postButton").style.display="block"
            document.getElementById("postButton").style.margin="0 auto"

            chatSocket.connectSocket("chat", chatSocket.chatHandler)
            presenceSocket.connectSocket("presence", presenceSocket.presenceHandler);
          
          // contentSocket.sendContentRequest(event)
          // postSocket.sendSubmittedPostsRequest() 
        }
      })

    }).catch((error)=>{
      console.log(error)

    })
    // console.log("logindata",logindata, "hi")
    // console.log( Object.keys(logindata).length)
    // console.log(JSON.stringify(logindata))
    postSocket.sendSubmittedPostsRequest()
  document.getElementById('login-form').reset()

 
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

  document.getElementById("postLogin").style.display="block"
  document.getElementById("postButton").style.display="none"
console.log("Logged out", response)
let presenceCont = document.getElementById("presencecontainer")
if (presenceCont.childElementCount != 0) {
  while (presenceCont.firstChild) {
  presenceCont.removeChild(presenceCont.lastChild);
}
}
})
postSocket.sendSubmittedPostsRequest()
let user= document.getElementById('welcome')
user.innerText = "Welcome"
}

function getCookieName(){
  let cookies = document.cookie.split(";")
  let lastCookieName = cookies[cookies.length -1].split("=")[0].replace(" ", '')
  // console.log("cookie",cookies, "length", cookies.length)
  return lastCookieName
  // console.log("h",lastCookieName)
}

let content = {
  postID :  "",
  types : "" 
}
function loadContent(type, postID){
  content.types = type
  content.postID = postID
  let stringified = JSON.stringify(content)
  fetch("/loadContent",{
    headers:{
      'Accept':'application/json',
      'Content-Type': 'application/json'
    },
    method: "POST",
    body: stringified
  })
  .then(response => response.json())
  .then(data =>{
    const c = JSON.parse(data)
    document.getElementById("content").innerHTML = c.body
  }).catch(error => console.log(error))
}

function Notifications(message) {
  fetch("/notification",{
    headers:{
      'Accept':'application/json',
      'Content-Type': 'application/json'
    },
    method: "POST",
    body: message
  })
  .then(response => response.json())
  .then((data)=>{
    //when confirmation is reciever add class to presence button that adds notification symbol 
    // let chat = JSON.parse(data)
    console.log(data.notification == true)
    if (data.notification == true){
      let button = document.querySelector(`.${data.sender.id}`)
      button.style.backgroundColor = "purple"
      console.log("notification added")
    }
  })
}

function reformatTime(date) {
  let hoursN = date.getHours()
  let minutesN = date.getMinutes() 
  let minutes
  if (minutesN<10) {
    minutes = "0" + minutesN.toString()
  } else {
    minutes = minutesN.toString()
  }
  let time = hoursN.toString() + ":" + minutes
  return time
}

