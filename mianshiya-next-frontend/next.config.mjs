/** @type {import('next').NextConfig} */
const nextConfig = {
    output: "standalone",
    typescript: {
        ignoreBuildErrors: true,
    },
    async rewrites() {
        return [
            {
                source: "/api/static/:path*",
                destination: "http://localhost:8101/api/static/:path*",
            },
        ];
    },
};

export default nextConfig;
