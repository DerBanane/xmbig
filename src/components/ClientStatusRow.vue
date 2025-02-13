<template>
  <div style="height: 30vh; margin: 0 auto">
    <canvas id="hashrateChart"></canvas>
    <br /><br />
    <canvas id="minerChart"></canvas>
  </div>
</template>

<script>
import { onMounted, ref } from 'vue';
import Chart from 'chart.js';
import $ from 'jquery';
import 'chartjs-plugin-colorschemes';


export default {
  props:{
    CHART_RELOAD_INTERVAL_IN_MS:{
      type:Number,
      required:true
    }
  },
  setup(props) {
    const hashrateChart = ref(null);
    const minerChart = ref(null);

    onMounted(() => {
      hashrateChart.value = createChart('hashrateChart', 'Hashrates', 'Hashrate');
      minerChart.value = createChart('minerChart', 'Miners', 'Miner', true);
      updateCharts();

      setInterval(() => {
        updateCharts();
      }, props.CHART_RELOAD_INTERVAL_IN_MS);
    });

    const createChart = (canvasId, title, labelString, integerValues = false) => {
      return new Chart(document.getElementById(canvasId).getContext('2d'), {
        type: 'line',
        data: { datasets: [] },
        options: {
          title: { display: true, text: title },
          tooltips: { mode: 'nearest', intersect: false },
          scales: {
            yAxes: [
              {
                id: canvasId.replace('Chart', ''), // Dynamische ID für y-Achse
                ticks: {
                  beginAtZero: true,
                  callback: integerValues
                    ? function (value) {
                        if (value % 1 === 0) {
                          return value;
                        }
                      }
                    : null,
                },
                scaleLabel: { display: true, labelString: labelString, fontSize: 16 },
              },
            ],
            xAxes: [
              {
                type: 'time',
                distribution: 'series',
                time: { unit: 'minute', displayFormats: { minute: 'HH:mm' } },
                ticks: { stepSize: 10, autoSkip: true, maxTicksLimit: 60 },
              },
            ],
          },
          elements: { point: { radius: 0 } },
          responsive: true,
          maintainAspectRatio: false,
          plugins: { colorschemes: { scheme: 'brewer.Paired12' } },
        },
      });
    };

    const updateChartData = (chart, results, dataLabel) => {
      var datasets = results['client_statistics'].map(function (entry) {
        return {
          label: entry.algo,
          type: 'line',
          data: entry.statistics.map(function (stat) {
            return { x: new Date(stat.timestamp), y: stat[dataLabel] };
          }),
          spanGraphs: false,
          fill: dataLabel === 'hashrate', // Nur für Hashrate-Chart füllen
        };
      });

      chart.data.datasets = datasets;
      chart.update();
    };

    const updateCharts = () => {
      $.ajax({
        url: '/admin/getClientStatistics',
        dataType: 'json',
      }).done((results) => {
          updateChartData(hashrateChart.value, results, 'hashrate');
          updateChartData(minerChart.value, results, 'miner');
        }
      );
    };

    return {
      hashrateChart,
      minerChart,
      createChart,
      updateChartData,
      updateCharts,
    };
  },
};
</script>