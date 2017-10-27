for (let i = 0; i < 1000; i++) {
    let intervalId = messageLoop()
    setTimeout(() => {
        clearInterval(intervalId)
    }, 10000)
}

function messageLoop() {
    let counter = 0;
    return setInterval(() => {
        counter++;
    }, 100)
}
