import request from '@/utils/request'

// 获取全部互信主机ip信息
export function getAllIpaddr() {
  return request({
    url: '/api/v1/ansible/inventory/hosts',
    method: 'get'
  })
}

// 获取inventory全部信息
export function getInventoryInfo() {
  return request({
    url: '/api/v1/ansible/inventory/groups',
    method: 'get'
  })
}

// 添加inventory
export function addInventoryInfo(data) {
  return request({
    url: '/api/v1/ansible/inventory/groups',
    method: 'post',
    data: data
  })
}
//  headers: { 'Content-Type': 'application/json' }

// 修改inventory
export function updateInventoryInfo(data) {
  return request({
    url: '/api/v1/ansible/inventory/groups',
    method: 'patch',
    data: data
  })
}

// 删除inventory
export function deleteInventoryInfo(group) {
  return request({
    url: '/api/v1/ansible/inventory/groups/' + group,
    method: 'delete'
  })
}
