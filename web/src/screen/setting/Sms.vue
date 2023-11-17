<script setup>
import {PlusOutlined, EditOutlined, DeleteOutlined} from "@ant-design/icons-vue";
import {computed, onMounted, ref} from "vue";
import AddSmsProfile from "../../components/AddSmsProfile.vue";
import useAxios from "../../support/axios.js";
import useMessage from "../../support/message.js";

const {httpGet, httpDelete} = useAxios()
const {errorMessage, successMessage} = useMessage()

const profiles = ref([])
const columns = ref([
    {
        title: 'Name',
        dataIndex: 'name',
        key: 'name',
    },
    {
        title: 'Provider',
        dataIndex: 'provider',
        key: 'provider',
    },
    {
        title: 'Access Key',
        dataIndex: 'ak',
        key: 'ak',
    },
    {
        title: 'From',
        dataIndex: 'from',
        key: 'from',
    },
    {
        title: 'Action',
        key: 'action',
    },
])
const profileFormRef = ref()
const loading = ref(false)
const keyword = ref()

const filteredProfiles = computed(() => {
    if(!keyword.value) return profiles.value
    return profiles.value.filter(p => {
        return p.name.indexOf(keyword.value) > -1 ||
            p?.provider?.indexOf(keyword.value) > -1 ||
            p?.ak?.indexOf(keyword.value) > -1 ||
            p?.from?.indexOf(keyword.value) > -1
    })
})

const fetchProfiles = () => {
    httpGet('/sms/profile').exec().then((data) => {
        profiles.value = data.payload
    }).catch(e => {
        errorMessage("Could not fetch profiles", e).show()
    })
}
const openAddSmsProfileForm = () => {
    profileFormRef.value.open()
}
const openEditProfileForm = (profile) => {
    profileFormRef.value.open(profile)
}
const handleDelete = (profile) => {
    loading.value = true
    httpDelete(`/sms/profile/${profile.id}`).exec().then(() => {
        successMessage("Profile deleted successfully").show()
        fetchProfiles()
    }).catch(e => {
        errorMessage("Could not delete profile", e).show()
    }).finally(() => {
        loading.value = false
    })
}

onMounted(() => {
    fetchProfiles()
})
</script>

<template>
    <a-row :gutter="[0, 12]">
        <a-col :span="12">
            <a-input-search
                class="template-search"
                v-model:value="keyword"
                placeholder="Enter keywords to start searching"
            />
        </a-col>
        <a-col :span="12" class="txt-rgt">
            <a-button type="primary" @click="openAddSmsProfileForm">
                <template #icon>
                    <PlusOutlined/>
                </template>
                New SMS profile
            </a-button>
        </a-col>
        <a-col :span="24">
            <a-table :columns="columns" :data-source="filteredProfiles" :loading="loading" :pagination="{defaultPageSize: 50}" size="small" bordered>
                <template #bodyCell="{ text, record, index, column }">
                    <template v-if="column.key === 'action'">
                        <a-space>
                            <a-button size="small" type="text" @click="openEditProfileForm(record)">
                                <template #icon>
                                    <EditOutlined color="blue"/>
                                </template>
                            </a-button>
                            <a-popconfirm
                                title="Are you sure delete this profile?"
                                ok-text="Yes"
                                cancel-text="No"
                                placement="bottom"
                                @confirm="handleDelete(record)"
                            >
                                <a-button size="small" type="text">
                                    <template #icon>
                                        <DeleteOutlined color="danger"/>
                                    </template>
                                </a-button>
                            </a-popconfirm>
                        </a-space>
                    </template>
                    <template v-if="column.key === 'provider'">
                        {{ $filters.smsProvider(text) }}
                    </template>
                </template>
            </a-table>
        </a-col>

        <AddSmsProfile ref="profileFormRef" @saved="fetchProfiles"></AddSmsProfile>
    </a-row>
</template>

<style scoped>

</style>
