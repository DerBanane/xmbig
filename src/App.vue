<template>
  <div style="width: 96%; margin: 0; background-color: DarkGray; padding: 10px; border: 1px solid black; border-radius: 10px;">
    <div id="updateNotificationBar"></div>
    <StatusBar :message="statusBarMessage" :type="statusBarType" />
    <SettingsPanel
      :HIDE_OFFLINE="HIDE_OFFLINE"
      :GROUP_BY_ALGO="GROUP_BY_ALGO"
      @hideOfflineChanged="handleHideOfflineChanged"
      @groupByAlgoChanged="handleGroupByAlgoChanged"
      @resetClientStatusList="handleResetClientStatusList"
    />
    <ClientStatusTable
      :TRESHOLD_IN_MS="TRESHOLD_IN_MS"
      :RELOAD_INTERVAL_IN_MS="RELOAD_INTERVAL_IN_MS"
      :HIDE_OFFLINE="HIDE_OFFLINE"
      :GROUP_BY_ALGO="GROUP_BY_ALGO"
      :currentServerTime="currentServerTime"
      @open-log-modal="openLogModal"
      @open-editor-modal="openEditorModal"
    />
    <DashboardCharts :CHART_RELOAD_INTERVAL_IN_MS="CHART_RELOAD_INTERVAL_IN_MS" />
    <MinerEditorModal
      v-if="showMinerEditorModal"
      :selectedClient="selectedClient"
      @close="closeMinerEditorModal"
      @config-saved="handleConfigSaved"
    />
    <MinerLogModal v-if="showMinerLogModal" :selectedClient="selectedClient" @close="closeMinerLogModal" />
  </div>
</template>

<script>
import ClientStatusTable from './components/ClientStatusTable.vue';
import DashboardCharts from './components/DashboardCharts.vue';
import SettingsPanel from './components/SettingsPanel.vue';
import StatusBar from './components/StatusBar.vue';
import MinerEditorModal from './components/MinerEditorModal.vue';
import MinerLogModal from './components/MinerLogModal.vue';
import { ref, onMounted } from 'vue';
import axios from 'axios';

export default {
  components: {
    ClientStatusTable,
    DashboardCharts,
    SettingsPanel,
    StatusBar,
    MinerEditorModal,
    MinerLogModal,
  },
  setup() {
    const TRESHOLD_IN_MS = 60 * 1000;
    const RELOAD_INTERVAL_IN_MS = 10 * 1000;
    const CHART_RELOAD_INTERVAL_IN_MS = 60 * 1000;

    const HIDE_OFFLINE = ref(true);
    const GROUP_BY_ALGO = ref(true);

    const selectedClient = ref(null);
    const showMinerEditorModal = ref(false);
    const showMinerLogModal = ref(false);

    const statusBarMessage = ref('');
    const statusBarType = ref('info');

    const currentServerTime = ref(Date.now()); // Initialisiere mit aktuellem Zeitstempel

    const setStatusBar = (message, type = 'info') => {
      statusBarMessage.value = message;
      statusBarType.value = type;
    };

    const handleHideOfflineChanged = (value) => {
      HIDE_OFFLINE.value = value;
    };

    const handleGroupByAlgoChanged = (value) => {
      GROUP_BY_ALGO.value = value;
    };

    const handleResetClientStatusList = async () => {
      try {
        await axios.post('/admin/resetClientStatusList', {});
        setStatusBar('Successfully sent the reset client status list request to the Server.  Now just wait for the next refresh.', 'success');
      } catch (error) {
        setStatusBar(`Failed to send the reset client status list request to the Server. \nError: ${error}`, 'danger');
      }
    };

    const openLogModal = (client) => {
      selectedClient.value = client;
      showMinerLogModal.value = true;
    };

    const closeMinerLogModal = () => {
      showMinerLogModal.value = false;
    };

    const openEditorModal = (client) => {
      selectedClient.value = client;
      showMinerEditorModal.value = true;
    };

    const closeMinerEditorModal = () => {
      showMinerEditorModal.value = false;
    };

    const handleConfigSaved = () => {
      setStatusBar('Successfully updated config for: ' + selectedClient.value.client_status.client_id, 'success');
    };

     onMounted(() => {
        // Aktualisiere die Serverzeit regelmäßig
        setInterval(() => {
            axios.get('/admin/getClientStatusList')
                .then(response => {
                    if (response.data && response.data.current_server_time) {
                        currentServerTime.value = response.data.current_server_time * 1000;
                    } else {
                        console.warn('current_server_time not found in response data');
                    }
                })
                .catch(error => {
                    console.error('Error fetching server time:', error);
                });
        }, RELOAD_INTERVAL_IN_MS); // Aktualisiere alle 10 Sekunden (RELOAD_INTERVAL_IN_MS)
    });

    return {
      TRESHOLD_IN_MS,
      RELOAD_INTERVAL_IN_MS,
      CHART_RELOAD_INTERVAL_IN_MS,
      HIDE_OFFLINE,
      GROUP_BY_ALGO,
      handleHideOfflineChanged,
      handleGroupByAlgoChanged,
      handleResetClientStatusList,
      selectedClient,
      showMinerEditorModal,
      showMinerLogModal,
      openLogModal,
      closeMinerLogModal,
      openEditorModal,
      closeMinerEditorModal,
      statusBarMessage,
      statusBarType,
      setStatusBar,
      handleConfigSaved,
      currentServerTime // Hinzugefügt
    };
  },
};
</script>