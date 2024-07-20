import { baseClient } from './base'

export async function getProviders(body) {
  const res = await baseClient.post('/providers', body)
  return res
}

export async function getProvidersExtraConfig(body, extraConfig) {
  const res = await baseClient.post('/providers', body, extraConfig)
  return res
}
