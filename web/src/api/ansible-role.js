import request from '@/utils/request'

// 添加role信息
export function addRoleInfo(role) {
  return request({
    url: '/api/v1/ansible/role/' + role,
    method: 'post'
  })
}

// 删除role
export function deleteRoleInfo(role) {
  return request({
    url: '/api/v1/ansible/role/' + role,
    method: 'delete'
  })
}

// 查询role信息
export function getRoleInfo(query, filter) {
  return request({
    url: '/api/v1/ansible/role',
    method: 'get',
    params: query, filter
  })
}
