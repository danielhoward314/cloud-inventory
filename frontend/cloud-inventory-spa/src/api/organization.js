import { baseClient } from './base'

export async function getOrganization(id) {
  const res = await baseClient.get(`/organizations/${id}`)
  return res
}
