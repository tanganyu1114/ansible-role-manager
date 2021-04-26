<template>
  <BasicLayout>
    <template #wrapper>
      <el-card class="box-card">
        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button
              v-permisaction="['ansible:inventory:add']"
              type="primary"
              icon="el-icon-upload2"
              size="small"
              @click="handleAdd"
            >上传</el-button>
          </el-col>
        </el-row>
        <!-- 表格信息 -->
        <el-table
          v-loading="loading"
          :data="roleList"
          style="width: 100%"
          border
          @selection-change="handleSelectionChange"
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
          <el-table-column
            label="标签"
            prop="tags"
            :filters="filterData"
            :filter-method="filterTag"
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
      <el-dialog :title="title" :visible.sync="open" width="800px">
        <el-upload
          class="upload-demo"
          drag
          action="https://jsonplaceholder.typicode.com/posts/"
          multiple
        >
          <i class="el-icon-upload" />
          <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
          <div slot="tip" class="el-upload__tip">只能上传zip/tgz文件</div>
        </el-upload>
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
      // 查询参数
      total: 0,
      queryParams: {
        pageIndex: 1,
        pageSize: 10
      }
    }
  },
  created() {
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
    },
    handleAdd() {
      this.title = '上传Role角色文件'
      this.open = true
    },
    handleUpdate(row) {
    },
    handleDelete() {
    },
    filterTag(value, row) {
      return row.tag === value
    },
    uploadRole() {
    },
    cancel() {
      this.open = false
    }
  }
}
</script>

<style scoped>

</style>
