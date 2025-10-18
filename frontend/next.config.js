/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  
  // API URL
  env: {
    API_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000/api',
  },

  // Images optimization
  images: {
    domains: [
      'picsum.photos',
      'api.freepik.com',
      'images.unsplash.com',
    ],
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**',
      },
    ],
  },

  // Rewrites for API proxy
  async rewrites() {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://backend:8000';
    return [
      {
        source: '/api/:path*',
        destination: `${apiUrl}/:path*`,
      },
    ];
  },
  
  output: 'standalone',
}

module.exports = nextConfig
