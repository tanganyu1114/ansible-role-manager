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
          v-loading="loading"
          :data="inventoryList"
          style="width: 100%"
          border
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" align="center" />
          <el-table-column type="index" label="序号" width="55" align="center" />
          <el-table-column label="组名" prop="groupName" align="center" width="200">
            <template slot-scope="scope">
              <el-tag type="success">
                {{ scope.row.groupName }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="IP地址" prop="ipAddrs">
            <template slot-scope="scope">
              <el-tag
                v-for="ip in scope.row.ipAddrs"
                :key="ip"
                class="tag-ip"
              >
                {{ ip }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="IP总数" prop="ipTotal" align="center" width="100" />
          <el-table-column label="操作" align="center" width="200" class-name="small-padding fixed-width">
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
                <el-input v-model.trim="form.groupName" placeholder="请输入组名" />
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="IP地址" prop="ipAddrs">
                <el-tag
                  v-for="ip in form.ipAddrs"
                  :key="ip"
                  v-model="form.ipAddrs"
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
      form: {
        groupName: '',
        ipAddrs: []
      },
      // 动态tag数据信息
      ipAddrs: [],
      inputVisible: false,
      inputIp: '',
      // 表格数据信息
      inventoryList: [],
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
        ],
        ipAddrs: [
          { required: true, message: 'IP地址不能为空', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.getList()
  },
  mounted() {
    // 挂载阶段
  },
  methods: {
    getList() {
      // todo: 把all和inventory信息合并放到后端统一返回，添加分页查询功能
      this.inventoryList = []
      getAllIpaddr().then(response => {
        if (response.code === 200) {
          const Ips = response.data
          const data = {}
          data.groupName = 'all'
          data.ipAddrs = Ips
          data.ipTotal = Ips.length
          this.inventoryList.unshift(data)
        } else {
          this.msgError(response.msg)
        }
      })
      getInventoryInfo().then(response => {
        if (response.code === 200) {
          const Gps = response.data
          for (const data of Object.values(Gps)) {
            data.ipTotal = data.ipAddrs.length
            this.inventoryList.push(data)
          }
        } else {
          this.msgError(response.msg)
        }
      })
      this.loading = false
    },
    resetForm() {
      this.form = {
        groupName: '',
        ipAddrs: []
      }
    },
    handleAdd() {
      this.resetForm()
      this.title = '新增inventory'
      this.open = true
    },
    handleUpdate(row) {
      this.resetForm()
      this.title = '修改inventory'
      this.form.groupName = row.groupName || this.ids[0].groupName
      this.form.ipAddrs = row.ipAddrs || this.ids[0].ipAddrs
      this.form.targetGroupName = this.form.groupName
      this.open = true
    },
    handleDelete(row) {
      const groupName = row.groupName || this.ids.map(item => item.groupName)
      this.$confirm(
        '是否确认删除组名为" ' + groupName + ' "的数据项?',
        '警告',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
        .then(function() {
          return deleteInventoryInfo(groupName)
        })
        .then(() => {
          this.getList()
          this.msgSuccess('删除成功')
        }).catch(function() {})
    },
    // dialog弹窗函数
    submitForm() {
      console.log('form:', this.form)
      if (this.form.targetGroupName === undefined) {
        addInventoryInfo(this.form).then(response => {
          if (response.code === 200) {
            this.msgSuccess('添加inventory成功')
            this.open = false
            this.getList()
          } else {
            this.msgError(response.msg)
          }
        })
      } else {
        updateInventoryInfo(this.form).then(response => {
          if (response.code === 200) {
            this.msgSuccess('修改inventory成功')
            this.open = false
            this.getList()
          } else {
            this.msgError(response.msg)
          }
        })
      }
    },
    cancel() {
      this.resetForm()
      this.open = false
    },
    // 动态tag函数信息
    handleClose(tag) {
      this.form.ipAddrs.splice(this.form.ipAddrs.indexOf(tag), 1)
    },

    showInput() {
      this.inputVisible = true
      this.$nextTick(_ => {
        this.$refs.saveTagInput.$refs.input.focus()
      })
    },

    handleInputConfirm() {
      const inIp = this.inputIp
      const ipSet = new Set(this.form.ipAddrs)
      if (this.validateIpaddr(inIp)) {
        ipSet.add(inIp)
        this.form.ipAddrs = [...ipSet]
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
    // 多选框选中数据
    handleSelectionChange(selection) {
      console.log(selection)
      // this.ids = selection.map(item => item.groupName)
      this.ids = selection
      this.single = selection.length !== 1
      this.multiple = !selection.length
    }
  }
}
</script>

<style scoped>
.el-tag + .el-tag {
  margin-left: 10px;
}
.tag-ip {
  margin-left: 10px;
  margin-bottom: 10px;
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
