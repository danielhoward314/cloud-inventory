export const parseJwt = (token) => {
  const base64Url = token.split('.')[1]

  // Replace URL-safe characters to base64 characters
  const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')

  // Decode the base64 string
  const jsonPayload = decodeURIComponent(
    atob(base64)
      .split('')
      .map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
      })
      .join('')
  )

  // Parse the JSON payload
  return JSON.parse(jsonPayload)
}
