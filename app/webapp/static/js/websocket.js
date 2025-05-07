// Create a new WebSocket connection
const ws = new WebSocket("wss://yourserver.example.com/path");

// Set the idle timeout period (10 minutes = 600000 ms)
const IDLE_TIMEOUT = 10 * 60 * 1000;
let idleTimer = null;

// Function to reset the idle timer
function resetIdleTimer() {
  // Clear the previous timer if it exists
  if (idleTimer) {
    clearTimeout(idleTimer);
  }
  // Start a new timer
  idleTimer = setTimeout(() => {
    console.log("No activity detected for 10 minutes. Closing WebSocket.");
    ws.close();
  }, IDLE_TIMEOUT);
}

// List of events that indicate user activity
const activityEvents = ['mousemove', 'keydown', 'click', 'scroll', 'touchstart'];
activityEvents.forEach(eventName => {
  document.addEventListener(eventName, resetIdleTimer, false);
});

// Initialize the idle timer when the page loads
resetIdleTimer();

// Optionally, log connection events for debugging
ws.onopen = () => {
  console.log("WebSocket connection opened.");
};

ws.onclose = (event) => {
  console.log("WebSocket connection closed.", event);
};

ws.onerror = (error) => {
  console.error("WebSocket encountered error:", error);
};

// Ensure the WebSocket is closed when the browser window is closed or refreshed
window.addEventListener("beforeunload", () => {
  ws.close();
});
