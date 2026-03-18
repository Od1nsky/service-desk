export const useApi = () => {
  const { $api } = useNuxtApp()
  return $api as import('axios').AxiosInstance
}
