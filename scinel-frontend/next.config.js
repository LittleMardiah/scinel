/** @type {import('next').NextConfig} */
const nextConfig = {
  typescript: { strictNullChecks: true },
  compress: true,
  poweredByHeader: false,
  reactStrictMode: true,
  swcMinify: true,
  images: { formats: ['image/avif', 'image/webp'] },
}
module.exports = nextConfig
