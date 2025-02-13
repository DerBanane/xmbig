<template>
  <div>
    <table id="clientStatusList" class="table table-striped table-bordered" cellspacing="0" width="100%">
      <thead>
        <tr>
          <th class="center" width="2%"><i class="fa fa-square-o" @click="toggleSelectAll"></i></th>
          <th>Worker Id</th>
          <th>Version</th>
          <!-- ... (Restliche Tabellenüberschriften) ... -->
        </tr>
      </thead>
      <tbody>
        <ClientStatusRow v-for="client in filteredClientStatusList" :key="client.client_status.client_id"
          :client="client" @logClicked="handleLogClicked" @editClicked="handleEditClicked" />
      </tbody>
      <tfoot>
        <tr>
          <th class="center" width="2%"><i class="fa fa-square-o" @click="toggleSelectAll"></i></th>
          <th class="left">Total:</th>
          <th></th>
          <!-- ... (Restliche Tabellen-Fußzeilen) ... -->
        </tr>
      </tfoot>
    </table>
  </div>
</template>

<script>
import { computed, onMounted, ref, watch } from 'vue';
import ClientStatusRow from './ClientStatusRow.vue';
import $ from 'jquery';
import 'datatables.net-bs4';
import 'datatables.net-buttons-bs4';
import 'datatables.net-select-bs4';
import 'datatables.net-buttons/js/buttons.colVis.js';
import 'datatables.net';
import 'datatables.net-buttons/js/buttons.html5.min.js';
import 'datatables.net-buttons/js/buttons.print.min.js';
import 'jszip';
import { version, algoAndPowVariantName, memory, cache, isOnline, uptime, clientInfo, clientStatus, round, laststatus } from '../utils/clientStatusFormatters.js'
import axios from 'axios';

export default {
  components: {
    ClientStatusRow,
  },
  props: {
    TRESHOLD_IN_MS: {
      type: Number,
      required: true,
    },
    RELOAD_INTERVAL_IN_MS: {
      type: Number,
      required: true,
    },
    HIDE_OFFLINE: {
      type: Boolean,
      required: true,
    },
    GROUP_BY_ALGO: {
      type: Boolean,
      required: true,
    },
    currentServerTime: {
      type: Number,
      required: true,
    }
  },
  emits: ['open-log-modal', 'open-editor-modal', 'setStatusBar'],
  setup(props, { emit }) {
    const clientStatusList = ref([]);
    const table = ref(null);
    const hideOffline = ref(props.HIDE_OFFLINE);
    const groupByAlgo = ref(props.GROUP_BY_ALGO);
    const lastStatusFn = computed(() => laststatus(props.TRESHOLD_IN_MS, props.currentServerTime));

    const filteredClientStatusList = computed(() => {
      return clientStatusList.value.filter((client) => isOnline(client.client_status.last_status_update * 1000, props.TRESHOLD_IN_MS, props.currentServerTime));
    });

    watch(
      () => props.HIDE_OFFLINE,
      (newValue) => {
        hideOffline.value = newValue;
        if (table.value) {
          table.value.draw(); // Redraw the table when the filter changes
        }
      }
    );

    watch(
      () => props.GROUP_BY_ALGO,
      (newValue) => {
        groupByAlgo.value = newValue;
        if (table.value) {
          if (newValue) {
            table.value.rowGroup().enable().draw();
            table.value.order.fixed({
              pre: [5, 'asc']
            }).draw();
          } else {
            table.value.rowGroup().disable().draw();
            table.value.order.fixed({
              pre: []
            }).draw();
          }
        }
      }
    );

    const tableOptions = {
      dom:
        "<'row'<'col-sm-12'B>><'row rowPadded'<'col-sm-9'l><'col-sm-3'f>><'row'<'col-sm-12't>><'row'<'col-sm-4'i><'col-sm-8'p>>",
      lengthMenu: [
        [10, 25, 50, 100, -1],
        [10, 25, 50, 100, "All"],
      ],
      pageLength: -1,
      bPpaginate: true,
      pagingType: "full_numbers",
      stateSave: true,
      ajax: {
        url: "/admin/getClientStatusList",
        dataSrc: (json) => {
          clientStatusList.value = json.client_status_list;
          return json.client_status_list;
        },
      },
      buttons: [
        'colvis',
        'copy',
        'excel',
        'csv',
        'print',
        {
          text: '<i class="fa fa-upload"> Sende miner config</i>',
          className: "btn-primary",
          enabled: false,
          action: function () {
            triggerAction('UPDATE_CONFIG');
          },
        },
        {
          text: '<i class="fa fa-download"> Empfange miner config</i>',
          className: "btn-info",
          enabled: false,
          action: function () {
            triggerAction('PUBLISH_CONFIG');
          },
        },
        {
          text: '<i class="fa fa-play"> Start</i>',
          className: "btn-success",
          enabled: false,
          action: function () {
            triggerAction('START');
          },
        },
        {
          text: '<i class="fa fa-pause"> Pause</i>',
          className: "btn-warning",
          enabled: false,
          action: function () {
            triggerAction('STOP');
          },
        },
        {
          text: '<i class="fa fa-repeat"> Restart</i>',
          className: "btn-info",
          enabled: false,
          action: function () {
            triggerAction("RESTART");
          },
        },
        {
          text: '<i class="fa fa-sign-out"> Stop</i>',
          className: "btn-danger",
          enabled: false,
          action: function () {
            triggerAction("SHUTDOWN");
          },
        },
        {
          text: '<i class="fa fa-refresh"> Reboot</i>',
          className: "btn-warning",
          enabled: false,
          action: function () {
            triggerAction("REBOOT");
          },
        },
        {
          text: '<i class="fa fa-rocket"> Ausfuehren</i>',
          className: "btn-success",
          enabled: false,
          action: function () {
            triggerAction("EXECUTE");
          },
        },
        {
          text: '<i class="fa fa-table"> Vorlage zuweisen</i>',
          className: "btn-info",
          enabled: false,
          action: function () {
            triggerAction("TEMPLATE");
          },
        },
        {
          text: '<i class="fa fa-edit"> Vorlagen Editor</i>',
          className: "btn-primary",
          enabled: true,
          action: function () {
            triggerAction("TEMPLATE_EDIT");
          },
        },
      ],
      orderFixed: [5, "asc"],
      rowGroup: {
        dataSrc: "client_status.current_algo_name",
      },
      columnDefs: [
        {
          targets: 1, // Spalte "Worker Id"
          createdCell: function (td, cellData, rowData, ) {
            $(td).attr('data-toggle', 'tooltip').attr('title', clientInfo(cellData, 'display', props.TRESHOLD_IN_MS, props.currentServerTime, rowData));
          }
        }
      ],
      columns: [
        {
          data: null,
          defaultContent: "",
          className: "select-checkbox",
          orderable: false,
        },
        { data: "client_status.client_id"},
        { data: "client_status.version", render: version },
        { data: "client_status.current_pool" },
        { data: "client_status.current_pool_user", visible: false },
        { data: "client_status.current_pool_pass", visible: false },
        { data: "client_status.current_status", render: (data, type, row) => clientStatus(data, props.TRESHOLD_IN_MS, props.currentServerTime, row) },
        { data: "client_status.current_algo_name", render: algoAndPowVariantName },
        { data: "client_status.cpu_brand", visible: false },
        { data: "client_status.external_ip", visible: false },
        { data: "client_status.hugepages_available", visible: false },
        { data: "client_status.hugepages_enabled", visible: false },
        { data: "client_status.cpu_is_x64", visible: false },
        { data: "client_status.cpu_has_aes", visible: false },
        { data: "client_status.cpu_is_vm", visible: false },
        { data: "client_status.hash_factor", className: "right", visible: false },
        { data: "client_status.total_pages", className: "right", visible: false },
        { data: "client_status.total_hugepages", className: "right", visible: false },
        { data: "client_status.free_memory", render: memory, className: "right", visible: false },
        { data: "client_status.total_memory", render: memory, className: "right", visible: false },
        { data: "client_status.current_threads", className: "right", visible: false },
        { data: "client_status.cpu_sockets", className: "right", visible: false },
        { data: "client_status.cpu_cores", className: "right", visible: false },
        { data: "client_status.cpu_threads", className: "right", visible: false },
        { data: "client_status.cpu_l2", render: cache, className: "right", visible: false },
        { data: "client_status.cpu_l3", render: cache, className: "right", visible: false },
        { data: "client_status.cpu_nodes", className: "right", visible: false },
        { data: "client_status.max_cpu_usage", className: "right", visible: false },
        { data: "client_status.hashrate_short", render: (data) => round(data), className: "right" },
        { data: "client_status.hashrate_medium", render: (data) => round(data), className: "right" },
        { data: "client_status.hashrate_long", render: (data) => round(data), className: "right" },
        { data: "client_status.hashrate_highest", render: (data) => round(data), className: "right" },
        { data: "client_status.hashes_total", className: "right" },
        { data: "client_status.avg_time", className: "right" },
        { data: "client_status.shares_good", className: "right" },
        { data: "client_status.shares_total", className: "right" },
        { data: "client_status.uptime", render: uptime, className: "right" },
         { data: "client_status.last_status_update", render: lastStatusFn.value },
        {
          data: null,
          defaultContent:
            "<td class='center-tab'><button type='button' id='LOG' class='btn btn-xs btn-info' data-toggle='tooltip' title='View miner log'><i class='fa fa-file-text-o'></i></button></td>",
          orderable: false,
          className: "center-tab",
        },
        {
          data: null,
          defaultContent:
            "<td class='center-tab'><button type='button' id='EDIT' class='btn btn-xs btn-primary' data-toggle='tooltip' title='Edit miner config'><i class='fa fa-edit'></i></button></td>",
          orderable: false,
          className: "center-tab",
        },
      ],
      rowId: "client_status.client_id",
      select: {
        style: "multi+shift",
      },
      order: [1, "asc"],
      footerCallback: function () {
        var api = this.api();

        var sumHashrateShort = 0;
        var sumHashrateMedium = 0;
        var sumHashrateLong = 0;
        var sumHashrateHighest = 0;
        var sumHashesTotal = 0;
        var avgTimeTotal = 0;
        var sumSharesGood = 0;
        var sumSharedTotal = 0;

        var colOffset = 28;

        sumHashrateShort = api
          .column(colOffset, { page: "current" })
          .data()
          .reduce(function (a, b) {
            return a + b;
          }, 0);

        sumHashrateMedium = api
          .column(colOffset + 1, { page: "current" })
          .data()
          .reduce(function (a, b) {
            return a + b;
          }, 0);

        sumHashrateLong = api
          .column(colOffset + 2, { page: "current" })
          .data()
          .reduce(function (a, b) {
            return a + b;
          }, 0);

        sumHashrateHighest = api
          .column(colOffset + 3, { page: "current" })
          .data()
          .reduce(function (a, b) {
            return a + b;
          }, 0);

        sumHashesTotal = api
          .column(colOffset + 4, { page: "current" })
          .data()
          .reduce(function (a, b) {
            return a + b;
          }, 0);

        avgTimeTotal = api
          .column(colOffset + 5, { page: "current" })
          .data()
          .reduce(function (a, b) {
            return a + b;
          }, 0) / api.column(26, { page: "current" }).data().length;

        sumSharesGood = api
          .column(colOffset + 6, { page: "current" })
          .data()
          .reduce(function (a, b) {
            return a + b;
          }, 0);

        sumSharedTotal = api
          .column(colOffset + 7, { page: "current" })
          .data()
          .reduce(function (a, b) {
            return a + b;
          }, 0);

        sumHashrateShort = round(sumHashrateShort);
        sumHashrateMedium = round(sumHashrateMedium);
        sumHashrateLong = round(sumHashrateLong);
        sumHashrateHighest = round(sumHashrateHighest);
        avgTimeTotal = round(avgTimeTotal);

        // update footer
        $(api.column(colOffset).footer()).html(sumHashrateShort);
        $(api.column(colOffset + 1).footer()).html(sumHashrateMedium);
        $(api.column(colOffset + 2).footer()).html(sumHashrateLong);
        $(api.column(colOffset + 3).footer()).html(sumHashrateHighest);
        $(api.column(colOffset + 4).footer()).html(sumHashesTotal);
        $(api.column(colOffset + 5).footer()).html(avgTimeTotal);
        $(api.column(colOffset + 6).footer()).html(sumSharesGood);
        $(api.column(colOffset + 7).footer()).html(sumSharedTotal);
      },
    };

    const triggerAction = (action, payload = null) => {
      table.value.rows({ selected: true }).eq(0).each(function (index) {
        var row = table.value.row(index);
        var data = row.data();
        const clientId = data.client_status.client_id;
        sendAction(action, clientId, payload);
      });
    };

    const updateButtonState = (selectedRows) => {
      for (let i = 5; i <= 14; i++) { // Buttons Start at index 5 now
        table.value.button(i).enable(selectedRows > 0);
      }
    }

    const toggleSelectAll = () => {
      if ($("#clientStatusList th i.fa-square-o").length) {
        table.value.rows().select();
        $("#clientStatusList th i")
          .removeClass("fa-square-o")
          .addClass("fa-check-square-o");
      } else {
        table.value.rows().deselect();
        $("#clientStatusList th i")
          .removeClass("fa-check-square-o")
          .addClass("fa-square-o");
      }
    };

    const sendAction = async (action, clientId, payload = null) => {
      try {
        await axios.post(`/admin/setClientCommand?clientId=${clientId}`, {
          control_command: { command: action, payload: payload },
        });
        emit('setStatusBar', `Successfully send ${action} to ${clientId} [payload='${payload}']`, 'success');
      } catch (error) {
        emit('setStatusBar', `Failed to send ${action} to ${clientId} [payload='${payload}'] \nError: ${error}`, 'danger');
      }
    };

    onMounted(() => {
      const $table = $('#clientStatusList');
      table.value = $table.DataTable(tableOptions);

      // Enable row grouping and apply the filter when the table is ready
      if (props.GROUP_BY_ALGO) {
        table.value.rowGroup().enable().draw();
        table.value.order.fixed({
          pre: [5, 'asc']
        }).draw();
      }

      $table.on('draw.dt', function () {
        // Re-initiate tooltips after each draw
        $('[data-toggle="tooltip"]').tooltip();
      });

      table.value.on('select deselect', function () {
        const selectedRows = table.value.rows({ selected: true }).count();
        updateButtonState(selectedRows);

        if ($("#clientStatusList th i.fa-square-o").length) {
          $("#clientStatusList th i")
            .removeClass("fa-square-o")
            .addClass("fa-check-square-o");
        }
        if (!selectedRows) {
          $("#clientStatusList th i")
            .removeClass("fa-check-square-o")
            .addClass("fa-square-o");
        }

      });

      //Apply row grouping if enabled
      if (props.GROUP_BY_ALGO) {
        table.value.rowGroup().enable().draw();
      }

      table.value.buttons().container().appendTo($table.closest('.dataTables_wrapper').find('.col-sm-12:first'));

      setInterval(function () {
        table.value.ajax.reload(null, false);
      }, props.RELOAD_INTERVAL_IN_MS);

      updateButtonState(0);

       // Initialisiere Tooltips
         $('[data-toggle="tooltip"]').tooltip();
    });

    const handleLogClicked = (client) => {
      emit('open-log-modal', client);
    };

    const handleEditClicked = (client) => {
      emit('open-editor-modal', client);
    };

    return {
      clientStatusList,
      filteredClientStatusList,
      table,
      toggleSelectAll,
      handleLogClicked,
      handleEditClicked,
      hideOffline,
      groupByAlgo,
      lastStatusFn
    };
  },
};
</script>

<style scoped>
.center {
  text-align: center;
}

.right {
  text-align: right;
}

.left {
  text-align: left;
}

.rowPadded {
  padding-top: 15px;
}

.online {
  color: green;
}

.offline {
  color: red;
}
</style>