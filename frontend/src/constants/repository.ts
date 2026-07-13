export const GITHUB_REPO = 'FLuoXue/sub2api'
export const GITHUB_URL = `https://github.com/${GITHUB_REPO}`
export const GITHUB_DOCS_URL = `${GITHUB_URL}/blob/main/docs`

// Release images are published to GHCR with lowercase repository names.
export const DOCKER_IMAGE = `ghcr.io/${GITHUB_REPO.toLowerCase()}`
