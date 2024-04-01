import service from '@/utils/request'

export const serverCheck = (data) => {
  return service({
    url: '/server/check',
    method: 'POST',
    data
  })
}
