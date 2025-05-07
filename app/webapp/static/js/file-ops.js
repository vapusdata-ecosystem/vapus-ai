import { FileServiceClient } from './generated/your_proto_package_grpc_web_pb';
import { FileRequest } from './generated/your_proto_package_pb';

const client = new FileServiceClient('http://localhost:8080', null, null); // Use the appropriate URL

const request = new FileRequest();
request.setFileName('your-file.txt');

const stream = client.streamFile(request, {});

stream.on('data', (response) => {
  const fileData = response.getData_asU8(); // Get data as Uint8Array

  // Process or display the data (e.g., append to an element or download)
  const blob = new Blob([fileData], { type: 'application/octet-stream' });
  const url = URL.createObjectURL(blob);

  // For example, trigger a download
  const link = document.createElement('a');
  link.href = url;
  link.download = 'your-file.txt';
  link.click();
  URL.revokeObjectURL(url);
});

stream.on('error', (err) => {
  console.error('Stream error:', err);
});

stream.on('end', () => {
  console.log('Stream ended');
});
