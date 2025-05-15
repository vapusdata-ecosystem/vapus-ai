/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    domains: ["i0.wp.com", "storage.googleapis.com"],
  },
  webpack: (config) => {
    config.module.rules.push({
      test: /\.css$/,
      use: ["style-loader", "css-loader"],
    });
    return config;
  },
};

export default nextConfig;
