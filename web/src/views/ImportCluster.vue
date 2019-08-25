<template>
    <div class="import-page">
        <p class="import-page__title">导入集群</p>
        <div class="import-container">
            <el-form :model="importForm" :rules="rules" ref="importForm" label-width="80px">

                <el-form-item label="业务线" prop="line" required>
                    <el-select v-model="importForm.line" placeholder="请选择业务线">
                        <el-option label="出借" value="lend"></el-option>
                        <el-option label="借款" value="loan"></el-option>
                        <el-option label="贷嘛" value="daima"></el-option>
                    </el-select>
                </el-form-item>

                <el-form-item label="Host" prop="host" required>
                    <el-input v-model="importForm.host"></el-input>
                </el-form-item>

                <el-form-item label="Port" prop="port" required>
                    <el-input v-model="importForm.port"></el-input>
                </el-form-item>


                <el-form-item label="密码" prop="password">
                    <el-input v-model="importForm.password"></el-input>
                </el-form-item>

                <el-form-item class="footer-item">
                    <el-button @click="checkForm('importForm')">验证</el-button>
                    <el-button type="primary" @click="submitForm('importForm')" :disabled="submitDisabled">导入
                    </el-button>
                </el-form-item>
            </el-form>
        </div>

    </div>
</template>

<script>
    import {getClusterNodesApi} from "@/http/api";

    export default {
        data() {
            const checkHost = (rule, value, callback) => {
                if (!value) {
                    callback(new Error('请输入host'))
                }
            };

            return {
                rules: {
                    line: [{
                        trigger: 'change'
                    }],
                    host: [{
                        // validator: checkHost,
                        trigger: 'change'
                    }]
                },
                importForm: {
                    line: "lend",
                    host: "10.211.55.12",
                    port: "8001",
                    password: '',
                },
                submitDisabled: false,
            }
        },
        created() {
        },
        methods: {
            checkForm(formName) {
                this.$refs[formName].validate((valid) => {
                    if (valid) {
                        this.onCheck();
                    } else {
                        return false
                    }
                });

            },
            async onCheck() {
                let params = JSON.parse(JSON.stringify(this.importForm));
                console.log(params);
                try {
                    await getClusterNodesApi(params);
                    this.$message.success("获取集群节点信息中，请等待")
                } catch ({error}) {
                    this.$message.error(`获取失败：${error}`)
                }
            },

            submitForm(formName) {
            }
        }
    }
</script>

<style lang="scss" scoped>
    @import '@/style/mixin.scss';

    $edit-icon-color: #1890ff;
    $green-color: #67C23A;

    .import-page {
        .el-input-group__append {
            width: 60px;
        }
    }

    .import-page__title {
        @include page-title-font;
        margin: 10px 0;
    }

    .import-container {
        padding: 20px 40px;
        margin-top: 20px;
        background: #fff;
        border-radius: 5px;

        .el-form {
            max-width: 600px;
        }
    }

    .type-tooltip {
        i {
            color: $green-color;
            font-size: 18px;
        }

        h6 {
            margin: 8px 0 5px 0;
        }

        p {
            margin: 3px 0;
        }

        .emphasis-text {
            color: #000;
        }
    }

    .type-icon {
        margin: 2px 10px;
    }

    .appid-select,
    .group-select {
        width: 100%;
    }

    .custom-spec {
        display: flex;

        .el-input {
            width: 150px;
            margin-right: 5px;

            .el-input-group__append {
                padding: 0 5px;
            }
        }
    }

    .footer-item {
        margin-top: 20px;
        display: inline;
        justify-content: flex-end;
    }
</style>