<template>
  <BasicLayout>
    <template #wrapper>
      <el-card class="box-card">
        <el-row :gutter="10" class="mb8" style="margin-bottom: 30px">
          <el-col :span="1.5">
            <el-button
              v-permisaction="['ansible:inventory:add']"
              type="primary"
              icon="el-icon-plus"
              size="small"
              @click="handleAdd"
            >新增</el-button>
          </el-col>
          <el-col :span="1.5">
            <el-button
              v-permisaction="['ansible:inventory:edit']"
              type="success"
              icon="el-icon-edit"
              size="small"
              :disabled="single"
              @click="handleUpdate"
            >修改</el-button>
          </el-col>
          <el-col :span="1.5">
            <el-button
              v-permisaction="['ansible:inventory:remove']"
              type="danger"
              icon="el-icon-delete"
              size="small"
              :disabled="multiple"
              @click="handleDelete"
            >删除</el-button>
          </el-col>
        </el-row>
        <!-- 表格信息 -->
        <el-table
          :data="tableData"
          style="width: 100%"
        >
          <el-table-column
            type="index"
            width="100"
          />
          <el-table-column
            label="组名"
            prop="groupName"
            width="200"
          />
          <el-table-column
            label="IP地址"
            width="500"
            prop="ipAddrs"
          />
          <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
            <template slot-scope="scope">
              <el-button
                v-permisaction="['ansible:inventory:edit']"
                size="mini"
                type="text"
                icon="el-icon-edit"
                @click="handleUpdate(scope.row)"
              >修改
              </el-button>
              <el-button
                v-permisaction="['ansible:inventory:remove']"
                size="mini"
                type="text"
                icon="el-icon-delete"
                @click="handleDelete(scope.row)"
              >删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <pagination
          v-show="total>10"
          :total="total"
          :page.sync="queryParams.pageIndex"
          :limit.sync="queryParams.pageSize"
          @pagination="getList"
        />
      </el-card>
      <!-- 添加或修改主机信息对话框 -->
      <el-dialog :title="title" :visible.sync="open" width="800px">
        <el-form ref="form" :model="form" :rules="rules" label-width="80px">
          <el-row>
            <el-col :span="24">
              <el-form-item label="组名" prop="groupName">
                <el-input v-model="form.groupName" placeholder="请输入组名" />
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="IP地址" prop="inputIp">
                <el-tag
                  v-for="ip in ipAddrs"
                  :key="ip"
                  closable
                  :disable-transitions="false"
                  @close="handleClose(ip)"
                >
                  {{ ip }}
                </el-tag>
                <el-input
                  v-if="inputVisible"
                  ref="saveTagInput"
                  v-model="inputIp"
                  class="input-ip"
                  size="medium"
                  @keyup.enter.native="$event.target.blur"
                  @blur="handleInputConfirm"
                />
                <el-button v-else class="btn-input-ip" @click="showInput">+ IP地址</el-button>
              </el-form-item>

            </el-col>
          </el-row>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-row>
            <el-col :span="16">
              <el-alert
                title="IP地址示例:10.1.1.10,10.1.1.[10:254]"
                type="info"
                :closable="false"
                show-icon
              />
            </el-col>
            <el-col :span="8">
              <el-button type="primary" @click="submitForm">确 定</el-button>
              <el-button @click="cancel">取 消</el-button>
            </el-col>
          </el-row>
        </div>
      </el-dialog>
    </template>
  </BasicLayout>
</template>

<script>
import { getInventoryInfo, getAllIpaddr, addInventoryInfo, updateInventoryInfo, deleteInventoryInfo } from '@/api/ansible-inventory'

export default {
  name: 'Inventory',
  data() {
    return {
      // 遮罩层
      loading: true,
      // 选中数组
      ids: [],
      // 弹出层标题
      title: '',
      // 非单个禁用
      single: true,
      // 非多个禁用
      multiple: true,
      // 是否显示弹出层
      open: false,
      form: {},
      // 动态tag数据信息
      ipAddrs: [],
      inputVisible: false,
      inputIp: '',
      // 表格数据信息
      tableData: [],
      // 查询参数
      total: 0,
      queryParams: {
        pageIndex: 1,
        pageSize: 10
      },
      // 表单校验
      rules: {
        groupName: [
          { required: true, message: '组名不能为空', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    (function getAllInfo() {
      const Ips = getAllIpaddr()
      // const Gps = getInventoryInfo()
      const data = {}
      data.groupName = 'all'
      data.ipAddrs = Ips
      this.tableData.push(data)
    }())
  },
  methods: {
    handleAdd() {
      this.title = '新增inventory'
      this.open = true
    },
    handleUpdate() {
      this.title = '修改inventory'
      this.open = true
    },
    handleDelete() {
      // TODO:
    },
    // dialog弹窗函数
    submitForm() {
      // TODO
      alert('form:' + this.form + 'ipaddr:' + this.form.ipAddrs + 'groupname:' + this.form.groupName)
    },
    cancel() {
      this.open = false
      // TODO
    },
    // 动态tag函数信息
    handleClose(tag) {
      this.ipAddrs.splice(this.ipAddrs.indexOf(tag), 1)
    },

    showInput() {
      this.inputVisible = true
      this.$nextTick(_ => {
        this.$refs.saveTagInput.$refs.input.focus()
      })
    },

    handleInputConfirm() {
      const inIp = this.inputIp
      if (this.validateIpaddr(inIp)) {
        this.ipAddrs.push(inIp)
      } else {
        this.$message.error('IP地址格式错误')
      }
      this.inputVisible = false
      this.inputIp = ''
    },
    // 验证ip地址是否合规
    validateIpaddr(ip) {
      const reg1 = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/
      const reg2 = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.\[(\d{1,2}|1\d\d|2[0-4]\d|25[0-5]):(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])]$/
      return reg1.test(ip) || reg2.test(ip)
    },

    // 表格函数
    getList() {
      this.loading = true
      /*      listSalaryTable(this.addDateRange(this.queryParams, this.dateRange)).then(response => {
          this.salarytableList = response.data.list
          this.total = response.data.count
          this.loading = false
        }
      )*/
    }
  }
}
</script>

<style scoped>
.el-tag + .el-tag {
  margin-left: 10px;
}
.btn-input-ip {
  margin-left: 10px;
  height: 32px;
  line-height: 30px;
  padding-top: 0;
  padding-bottom: 0;
}
.input-ip {
  width: 150px;
  margin-left: 10px;
  vertical-align: bottom;
}
</style>
