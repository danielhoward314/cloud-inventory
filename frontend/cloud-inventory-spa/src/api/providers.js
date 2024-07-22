import { baseClient } from './base'

export async function getProviders(organizationId) {
  const res = await baseClient.get(`/providers/${organizationId}`)
  return res
}
