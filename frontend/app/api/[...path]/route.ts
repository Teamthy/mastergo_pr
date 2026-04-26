import { NextRequest, NextResponse } from 'next/server'

const BACKEND = process.env.BACKEND_URL || "http://localhost:8080"

async function handler(req: NextRequest) {
    const url = new URL(req.url)

    const path = url.pathname.replace(/^\/api/, "")
    const target = `${BACKEND}${path}${url.search}`

    const res = await fetch(target, {
        method: req.method,
        headers: {
            "content-type": "application/json",
            cookie: req.headers.get("cookie") || "",
        },
        body:
            req.method !== "GET" && req.method !== "HEAD"
                ? await req.text()
                : undefined,
    })

    return new NextResponse(res.body, {
        status: res.status,
        headers: {
            "content-type": "application/json",
        },
    })
}

export { handler as GET, handler as POST, handler as DELETE, handler as PUT }