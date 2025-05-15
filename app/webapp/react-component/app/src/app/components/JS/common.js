// /**
//  * Formats a string to title case with optional text removal
//  * @param tr - Text to remove from the input string (can be empty)
//  * @param cc - Input string to format
//  * @returns The formatted string in title case
//  */

// export function strTitle(tr, cc) {
//   cc = cc.replace(/_/g, " ");
//   cc = cc.toUpperCase();
//   if (tr !== "") {
//     cc = cc.replace(new RegExp(tr, "gi"), "");
//   }
//   cc = cc.toLowerCase().replace(/(?:^|\s)\S/g, function (match) {
//     return match.toUpperCase();
//   });

//   return cc;
// }

/**
 * START THE STRING TO TITLE CASE FUNCTION
 * @param {string} input
 * @returns {string}
 */
export function strTitle(input) {
  if (!input) return "";
  let result = input.replace(/_/g, " ");
  return result
    .toLowerCase()
    .split(" ")
    .map((word) => {
      return word.charAt(0).toUpperCase() + word.slice(1);
    })
    .join(" ");
}
// END THE STRING TO TITLE CASE FUNCTION

// function for character Limit
/**
 * @param {string} input
 * @param {number} limit
 * @returns {string}
 */
export function wordLimit(input, limit) {
  if (!input) return "";
  if (!limit || limit <= 0) return input;
  if (input.length <= limit) {
    return input;
  }
  return input.substring(0, limit) + "...";
}

// End

// ToolTips function start
(function initTooltips() {
  if (typeof window === "undefined" || typeof document === "undefined") return;

  function addTooltipToElement(element) {
    if (!element || element._hasTooltip) return;
    element._hasTooltip = true;
    // Add event listeners
    element.addEventListener("mouseenter", function () {
      const tooltip = document.createElement("div");
      tooltip.style.position = "absolute";
      tooltip.style.backgroundColor = "#27272a";
      tooltip.style.color = "#fff";
      tooltip.style.padding = "5px 10px";
      tooltip.style.borderRadius = "4px";
      tooltip.style.fontSize = "14px";
      tooltip.style.zIndex = "1000";
      tooltip.style.pointerEvents = "none";
      // Set tooltip text to element's text content
      tooltip.textContent = this.getAttribute("tooltip") || this.textContent;
      document.body.appendChild(tooltip);
      const rect = this.getBoundingClientRect();
      tooltip.style.left = rect.left + "px";
      tooltip.style.top = rect.top - tooltip.offsetHeight - 5 + "px";

      this._tooltip = tooltip;
    });

    element.addEventListener("mouseleave", function () {
      if (this._tooltip && this._tooltip.parentNode) {
        document.body.removeChild(this._tooltip);
        this._tooltip = null;
      }
    });
  }

  // Function to initialize tooltips in the document
  function initializeTooltips() {
    const elementsWithTooltip = document.querySelectorAll(".tooltip");
    elementsWithTooltip.forEach(addTooltipToElement);
  }
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initializeTooltips);
  } else {
    initializeTooltips();
  }

  // Set up a mutation observer to handle dynamically added elements
  const observer = new MutationObserver(function (mutations) {
    mutations.forEach(function (mutation) {
      if (mutation.addedNodes.length) {
        mutation.addedNodes.forEach(function (node) {
          if (node.nodeType === 1) {
            if (node.classList && node.classList.contains("tooltip")) {
              addTooltipToElement(node);
            }
            if (node.querySelectorAll) {
              const tooltipElements = node.querySelectorAll(".tooltip");
              tooltipElements.forEach(addTooltipToElement);
            }
          }
        });
      }
    });
  });
  observer.observe(document.body, { childList: true, subtree: true });
  window.addEventListener("unload", function () {
    observer.disconnect();
  });
})();

//  TOOLTIP FUNCTION END HERE
