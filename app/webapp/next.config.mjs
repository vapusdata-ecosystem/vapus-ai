/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    domains: ["i0.wp.com", "storage.googleapis.com"],
  },
  webpack: (config) => {
    config.module.rules.push({
      test: /\.css$/,
      exclude: /codemirror/,
      use: ["style-loader", "css-loader"],
    });
    config.module.rules.push({
      test: /codemirror.*\.css$/,
      use: [
        'style-loader',
        {
          loader: 'css-loader',
          options: {
            url: false, // Disable URL handling for CodeMirror
            import: false, // Disable @import handling
          },
        },
      ],
    });
    
    return config;
  },
  reactStrictMode: false,
};

export default nextConfig;
