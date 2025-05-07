module.exports = {
    content: [
      "templates/*.html", // Adjust this to match your template files
      // "./templates/*.tmpl", // For Go's templating files
      "static/js/*.js",   // If you're using JavaScript
    ],
    theme: {
      extend: {},
    },
    plugins: [
      require('tailwind-scrollbar')({ nocompatible: true }),
    ],
  };
  