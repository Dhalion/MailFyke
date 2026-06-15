import { createAPIClient } from "openapi-fetch"
import type { paths } from "../types/api"

export const api = createAPIClient<paths>({
  baseUrl: "/api",
})
