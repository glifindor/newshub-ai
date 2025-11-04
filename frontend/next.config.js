/** @type {import('next').NextConfig} */
const publicApiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000/api'
const internalApiUrl = process.env.INTERNAL_API_URL || 'http://backend:8000/api'

const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  
  // API URL
  env: {
    API_URL: publicApiUrl,
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
    const normalizedInternalUrl = internalApiUrl.replace(/\/$/, '')
    return [
      {
        source: '/api/:path*',
        destination: `${normalizedInternalUrl}/:path*`,
      },
    ];
  },
  
  output: 'standalone',
}

module.exports = nextConfig
