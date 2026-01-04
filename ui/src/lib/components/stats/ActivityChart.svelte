<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import Chart from "chart.js/auto";
  import type { TimelineDataPoint } from "$lib/api";

  interface Props {
    data: TimelineDataPoint[];
  }

  let { data }: Props = $props();

  let canvas: HTMLCanvasElement;
  let chart: Chart | null = null;

  onMount(() => {
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    chart = new Chart(ctx, {
      type: "bar",
      data: {
        labels: data.map((d) => d.month),
        datasets: [
          {
            label: "Activities",
            data: data.map((d) => d.activities),
            backgroundColor: "hsl(221.2 83.2% 53.3% / 0.8)",
            borderRadius: 4,
          },
          {
            label: "Notes",
            data: data.map((d) => d.notes),
            backgroundColor: "hsl(215.4 16.3% 46.9% / 0.5)",
            borderRadius: 4,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            position: "top",
            labels: {
              usePointStyle: true,
              padding: 20,
            },
          },
        },
        scales: {
          x: {
            grid: {
              display: false,
            },
          },
          y: {
            beginAtZero: true,
            ticks: {
              stepSize: 1,
            },
            grid: {
              color: "hsl(215.4 16.3% 46.9% / 0.1)",
            },
          },
        },
      },
    });
  });

  onDestroy(() => {
    chart?.destroy();
  });

  $effect(() => {
    if (chart && data) {
      chart.data.labels = data.map((d) => d.month);
      chart.data.datasets[0].data = data.map((d) => d.activities);
      chart.data.datasets[1].data = data.map((d) => d.notes);
      chart.update();
    }
  });
</script>

<div class="h-64 w-full">
  <canvas bind:this={canvas}></canvas>
</div>
