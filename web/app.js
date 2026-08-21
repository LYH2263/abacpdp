async function refresh() {
  const [stats, policy] = await Promise.all([
    fetch('/api/stats').then(r => r.json()),
    fetch('/api/policy').then(r => r.json()),
  ]);
  document.getElementById('stats').textContent = JSON.stringify(stats, null, 2);
  document.getElementById('policy').textContent = JSON.stringify(policy, null, 2);
}

document.getElementById('eval').onclick = async () => {
  const bag = {
    subject: JSON.parse(document.getElementById('subject').value || '{}'),
    resource: JSON.parse(document.getElementById('resource').value || '{}'),
    action: JSON.parse(document.getElementById('action').value || '{}'),
    environment: JSON.parse(document.getElementById('env').value || '{}'),
  };
  const res = await fetch('/api/evaluate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(bag),
  });
  document.getElementById('result').textContent = await res.text();
  refresh();
};

refresh();
setInterval(refresh, 5000);
