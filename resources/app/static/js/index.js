let index = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            index.listen();
            index.modalEvent();
            // Event Manager
            index.startEvent();
            index.exitEvent();
        })
    },
    listen: function() {
        astilectron.onMessage(function(message) {
            switch (message.name) {
                case "downloadMissing":
                    webview = document.getElementById("mainScreen");
                    args = message.payload.split(";")
                    argString = `displayDownload("${args[0]}","${args[1]}");`;
                    webview.executeJavaScript(argString);
                    break;
                case "downloadFinished":
                    webview = document.getElementById("mainScreen");
                    webview.executeJavaScript(`stopDisplayDownload()`);
                    break;
                case "setupFinished":
                    var startButton = document.getElementById("StartButton");
                    startButton.innerText = "Start Wallet";
                    index.startWalletEvent();
                    break;
                case "check.out.menu":
                    asticode.notifier.info(message.payload);
                    break;
            }
        });
    },
    exitEvent: function(){
        var exitButton = document.getElementById("ExitButton");
        // When the user clicks on the button, open the modal
        exitButton.onclick = function() {
            let message = {"name": "exitEvent"};
            astilectron.sendMessage(message);
        }
    },
    startEvent: function(){
        var startButton = document.getElementById("StartButton");
        startButton.onclick = function (){
            index.startWalletEvent();
        };
    },

    modalEvent: function(){
        // Get the modal
        var modal = document.getElementById('aboutDialog');
        // Get the button that opens the modal
        var btn = document.getElementById("openAbout");
        // Get the <span> element that closes the modal
        var span = document.getElementsByClassName("close")[0];
        // When the user clicks on the button, open the modal
        btn.onclick = function() {
            modal.style.display = "block";
        }
        // When the user clicks on <span> (x), close the modal
        span.onclick = function() {
            modal.style.display = "none";
        }

        // When the user clicks anywhere outside of the modal, close it
        window.onclick = function(event) {
            if (event.target == modal) {
                modal.style.display = "none";
            }
        }
    },
    startWalletEvent: function(){
        var mainView = document.getElementById("mainScreen");
        var startButton = document.getElementById("StartButton");
        if (startButton.innerText == "Start Wallet"){
            let message = {"name": "startEvent"};
            astilectron.sendMessage(message);
            mainView.style.opacity = 0;
            setTimeout(function () {
                mainView.src = "static/html/loadscreen.html";
                setTimeout(function () {
                    mainView.style.opacity = 1;
                }, 500);
            }, 500);
            startButton.innerText = "Stop Wallet";
            startButton.className = "btn btn-outline-warning";
        }
        else{
            let message = {"name": "stopEvent"};
            astilectron.sendMessage(message);
            mainView.style.opacity = 0;
            setTimeout(function () {
                mainView.src = "static/html/stoppingscreen.html";
                setTimeout(function () {
                    mainView.style.opacity = 1;
                }, 500);
            }, 500);
            startButton.disabled = true;
        }
    },
};