# itemsim-server

## Updating KMS Gear Data

1. Export gear data and icon images using `WzJson` .
2. Replace `src/data/gear.json`, `src/data/gear-origin.json`, `src/data/gear-raw-origin.json`.
3. Run `pnpm wrangler deploy`.
4. Cd to `gearicon` and run `aws s3 sync . s3://{bucket-name}/gears/icon`
5. Cd to `gearrawicon` and run `aws s3 sync . s3://{bucket-name}/gears/iconRaw`
