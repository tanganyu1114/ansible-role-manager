<template>
  <BasicLayout>
    <template #wrapper>
      <el-card class="box-card">
        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button
              v-permisaction="['ansible:role:add']"
              type="primary"
              icon="el-icon-upload2"
              size="small"
              @click="handleAdd"
            >上传</el-button>
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
          :data="roleList"
          style="width: 100%"
          border
          @selection-change="handleSelectionChange"
          @filter-change="handleFilterChange"
        >
          <el-table-column type="selection" width="55" align="center" />
          <el-table-column type="index" label="序号" width="55" align="center" />
          <el-table-column label="角色名" prop="roleName" align="center" width="200">
            <template slot-scope="scope">
              <el-tag type="success">
                {{ scope.row.roleName }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="描述" prop="roleComment" />
          <el-table-column
            label="标签"
            prop="tags"
            :filters="filterData"
            :filter-method="filterTag"
            column-key="tags"
            filter-placement="bottom-end"
          >
            <template slot-scope="scope">
              <el-tag
                v-for="tag in scope.row.tags"
                :key="tag"
              >
                {{ tag }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" align="center" width="200" class-name="small-padding fixed-width">
            <template slot-scope="scope">
              <el-button
                v-permisaction="['ansible:role:edit']"
                size="mini"
                type="text"
                icon="el-icon-edit"
                @click="handleUpdate(scope.row)"
              >修改
              </el-button>
              <el-button
                v-permisaction="['ansible:role:remove']"
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
          :page.sync="queryParams.page"
          :limit.sync="queryParams.limit"
          @pagination="getList"
        />
      </el-card>
      <el-dialog :title="title" :visible.sync="open" width="800px">
        <el-row>
          <el-col :span="12">
            <el-upload
              ref="upload"
              class="upload-role"
              accept=".zip, .tgz"
              :limit="1"
              :auto-upload="false"
              :headers="upload.headers"
              :on-change="handleOnChange"
              :before-upload="handleBeforeUpload"
              :on-success="handleOnSuccess"
              :on-error="handleOnError"
              :data="upload.roleName"
              :action="upload.uploadUrl + form.uploadName"
              drag
            >
              <i class="el-icon-upload" />
              <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
              <div slot="tip" class="el-upload__tip">只能上传zip/tgz文件</div>
            </el-upload>
          </el-col>
          <el-col :span="12">
            <el-form ref="form" :model="form" :rules="rules">
              <el-row :gutter="20">
                <el-col :span="20">
                  <el-form-item label="角色名称" prop="uploadName">
                    <el-input v-model.trim="form.uploadName" placeholder="请输入角色名称" clearable />
                  </el-form-item>
                </el-col>
              </el-row>
            </el-form>
          </el-col>
        </el-row>
        <div slot="footer" class="dialog-footer">
          <el-row>
            <el-col>
              <el-button type="primary" @click="uploadRole">确 定</el-button>
              <el-button @click="cancel">取 消</el-button>
            </el-col>
          </el-row>
        </div>
      </el-dialog>
    </template>
  </BasicLayout>
</template>

<script>
import { getToken } from '@/utils/auth'

export default {
  name: 'Role',
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
      // role表格数据
      roleList: [],
      // tag筛选数据
      filterData: [],
      // upload
      upload: {
        // 上传role名字
        roleName: {},
        // 上传文件时头信息
        headers: { Authorization: 'Bearer ' + getToken() },
        // 上传地址
        uploadUrl: process.env.VUE_APP_BASE_API + '/api/v1/ansible/roles/'
      },
      form: {
        uploadName: ''
      },
      // 查询参数
      total: 0,
      queryParams: {
        page: 1,
        limit: 10,
        tags: []
      },
      // 表单校验
      rules: {
        uploadName: [
          { required: true, message: '角色名称不能为空', trigger: 'blur' },
          { pattern: /^(\w){4,20}$/, message: '只能输入4-20个字母、数字、下划线' }
        ]
      }
    }
  },
  created() {
    this.getList()
  },
  mounted() {
  },
  methods: {
    // 多选框选中数据
    handleSelectionChange(selection) {
      console.log(selection)
      // this.ids = selection.map(item => item.groupName)
      this.ids = selection
      this.single = selection.length !== 1
      this.multiple = !selection.length
    },
    getList() {
      this.loading = false
    },
    handleAdd() {
      this.title = '上传Role角色文件'
      this.open = true
    },
    handleUpdate(row) {
    },
    handleDelete() {
    },
    handleFilterChange(filter) {
      console.log(filter.tags)
      // TODO: 这里去请求后端role数据信息
      // queryParams.tags = filter.tags
      // getRoleInfo(queryParams)
    },
    filterTag(value, row) {
      return row.tag === value
    },
    uploadRole() {
      this.upload.roleName = { 'role': this.form.uploadName }
      // 手动上传
      this.$refs.upload.submit()
    },
    cancel() {
      this.open = false
    },
    handleOnChange(file) {
      this.form.uploadName = this.form.uploadName || file.name.split('.')[0]
    },
    handleBeforeUpload(file) {
      const fileArr = file.name.split('.')
      const types = ['zip', 'tgz']
      if (types.includes(fileArr[fileArr.length - 1])) {
        return true
      } else {
        this.msgError('文件格式错误!')
      }
    },
    handleOnSuccess(response, file, fileList) {
      console.log('response:', response)
      console.log('file', file)
      console.log('filelist', fileList)
    },
    handleOnError(err, file, fileList) {
      console.log('response:', err)
      console.log('file', file)
      console.log('filelist', fileList)
    }
  }
}
</script>

<style scoped>

</style>
