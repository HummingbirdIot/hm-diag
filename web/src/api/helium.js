export async function fetchHeliumHeight() {
  const api = `/inner/api/v1/proxy/heliumApi?path=/v1/blocks/height&t=${Date.now()}`;
  console.log(api);
  return new Promise((resolve, reject) => {
    fetch(api)
      .then((r) => r.json())
      .then((r) => {
        if (r.code === 200) {
          const d = JSON.parse(r.data);
          resolve(d.data.height);
        } else {
          reject(`fetch helium block height error: ${r.message}`);
        }
      })
      .catch((err) => {
        reject(err);
      });
  });
}
