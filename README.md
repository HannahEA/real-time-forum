TRELLO BOARD: https://trello.com/b/YQclPFFE/real-time-forum

PROJECT: https://learn.01founders.co/git/root/public/src/branch/master/subjects/real-time-forum

AUDIT: https://learn.01founders.co/git/root/public/src/branch/master/subjects/real-time-forum/audit


USE DEV=TRUE FOR DATABASE INIT ----- DEV=true go run .

RESOURCES:
    - https://pkg.go.dev/github.com/gorilla/websocket
    - https://medium.com/@bootcampmillionaire/what-i-learned-about-websockets-by-building-a-real-time-chat-application-using-socket-io-3d9e163e504
    - https://javascript.info/websocket
    - https://medium.com/@antoharyanto/make-simple-chat-application-using-golang-websocket-and-vanilla-js-f600e8020961
        - https://github.com/nnttoo/Simple-Chat-Apps-Goalang
    - https://www.figma.com/community/file/1114929434468857549

WebSockets in backend/frontend
    - https://developer.mozilla.org/en-US/docs/Web/API/WebSocket#examples
    - https://tutorialedge.net/golang/go-websocket-tutorial/
        - https://github.com/TutorialEdge/go-websockets-tutorial
    - websockets depend on two hosts having a connection, so they run on top of a TCP layer and modify the TCP layer so that the client and server agree for the socket to stay open
        - https://en.wikipedia.org/wiki/Transmission_Control_Protocol
        - https://en.wikipedia.org/wiki/Connection-oriented_communication
    - websocket handshake: how the data being exchanged should be interpreted by both client and server 
    - user does not have to refresh the page to see

- one HTML file only: https://en.wikipedia.org/wiki/Single-page_application

- Registration and login
    - form data must include: nickname, age, gender, first and last name, email, password
    - login using either nickname or email with password
    - logout from any page
    - add authentication for users/guests

- Creation of posts and comments
    - posts will have categories
    - does not have to be live

- Private messaging
    - online/offline section
        - organised by last message sent
            - if user is new, then contacts are ordered alphabetical
        - send private messages to online users
        - must be visible at all times
    - when user clicked on, it loads previous messages with that user
        - reloads the last 10 messages and when user scrolls up, 10 more messages should be shown (without spamming the scroll event)
            - look at [throttle](https://css-tricks.com/debouncing-throttling-explained-examples/#throttle)/debounce for not spamming the scroll event to see messages
    - message format
        - date and nickname shown

Bonus
- user profiles
- send images through messaging
- code is synchronicity (promises and go routines/channels) to increase performance

TODO:

https://medium.com/@antoharyanto/make-simple-chat-application-using-golang-websocket-and-vanilla-js-f600e8020961

- Logout/Auth:
    - finish login/register page, link up to the database
    - login asks for either nickname OR email with password to login, so ensure the funcs allow for this
    - add logout button on the header into the index.html, as it needs to be on the page at all times
- Chat:
    - users organised by online/offline in the presence list, online at the top, with last contacted at the top
    - add date to each chat message
    - add notifications for DMing
    - display 10 chat messages at a time, scroll to see more (throttle)
- Post:
    - fix bug with submitted posts not always rendering when you click the home button






