// document.addEventListener('DOMContentLoaded', function() {
//     const domainDropdownButton = document.getElementById('domainDropdownButton');
//     const domainItems = document.getElementById('domainItems');

//     domainDropdownButton.addEventListener('click', function(event) {
//         domainItems.classList.toggle('hidden');
//         event.stopPropagation(); // Prevent the click event from propagating to the document
//     });

//     // Close dropdown if clicked outside
//     document.addEventListener('click', function(event) {
//         if (!domainDropdownButton.contains(event.target) && !domainItems.contains(event.target)) {
//             domainItems.classList.add('hidden');
//         }
//     });

//     // Prevent the click event from propagating to the document
//     domainItems.addEventListener('click', function(event) {
//         event.stopPropagation();
//     });
// });
// document.addEventListener('DOMContentLoaded', function() {
//     document.getElementById('domainDropdownButton').addEventListener('click', function(event) {
//         console.log('clicked');
//         console.log
//         event.stopPropagation(); // Prevent the click event from bubbling up to the document
//         var domainItems = document.getElementById('domainItems');
//         if (domainItems.classList.contains('hidden')) {
//             domainItems.classList.remove('hidden');
//         } else {
//             domainItems.classList.add('hidden');
//         }
//     });

//     // Close the dropdown if clicked outside
//     document.addEventListener('click', function(event) {
//         var domainItems = document.getElementById('domainItems');
//         if (!document.getElementById('domainDropdownButton').contains(event.target)) {
//             domainItems.classList.add('hidden');
//         }
//     });
// });