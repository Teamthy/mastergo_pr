export async function fetchJSON(url: string, options?: RequestInit) {
    const res = await fetch(url, options)

    const text = await res.text()

    if (!text) return null

    try {
        return JSON.parse(text)
    } catch (e) {
        console.error("Invalid JSON:", text)
        throw e
    }
}