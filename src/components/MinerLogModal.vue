<template>
  <div class="modal fade" id="minerLog" tabindex="-1" role="dialog" aria-labelledby="minerLogLabel" aria-hidden="true">
    <div class="modal-dialog modal-lg" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="minerLogLabel">Miner Log</h5>
          <button type="button" class="close" data-dismiss="modal" aria-label="Close" @click="close">
            <span aria-hidden="true">×</span>
          </button>
        </div>
        <div class="modal-body">
          <div class="form-group" v-if="selectedClient" :data-client-id="selectedClient.client_status.client_id">
            <label for="log">Log of: {{ selectedClient.client_status.client_id }} ({{ selectedClient.client_status.external_ip }})</label>
            <textarea class="form-control" rows="20" id="log" v-model="clientLog" readonly></textarea>
          </div>
          <div v-else>
            Kein Client ausgewählt.
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-success" @click="refreshLog">Refresh</button>
          <button type="button" class="btn btn-secondary" data-dismiss="modal" @click="close">Close</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import $ from 'jquery';

export default {
  props: {
    selectedClient: {
      type: Object,
      default: null,
    },
  },
  emits: ['close', 'setStatusBar'],
  setup(props, { emit }) {
    const clientLog = ref('');

    const refreshLog = async () => {
      try {
        const response = await axios.get(`/admin/getClientLog?clientId=${props.selectedClient.client_status.client_id}`);
        clientLog.value = response.data.client_log;
      } catch (error) {
        console.error('Error fetching client log:', error);
        emit('setStatusBar', `Unable to fetch client log for ${props.selectedClient.client_status.client_id}. Please make sure it is enabled on the miner!`, 'danger');
      }
    };

    const close = () => {
      emit('close');
      $('#minerLog').modal('hide');
    };

    onMounted(() => {
      refreshLog();
      $('#minerLog').modal('show');
    });

    return {
      clientLog,
      refreshLog,
      close,
    };
  },
};
</script>