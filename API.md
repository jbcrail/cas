## API

**Get:**

    $ curl -X GET http://localhost:8000/:id

Get a specific item. The `:id` represents the SHA-1 of the requested content.

Possible HTTP status codes:

<table>
  <tr>
    <th>Status Code</th><th>Description</th>
  </tr>
  <tr>
    <td>200</td><td>Content exists, valid SHA-1</td>
  </tr>
  <tr>
    <td>404</td><td>Non-existent data requested</td>
  </tr>
  <tr>
    <td>500</td><td>Content is corrupted</td>
  </tr>
</table>

**Post:**

    $ curl -X POST --data-binary :filename http://localhost:5353/

Save the POSTed content to storage. The reply will contain the SHA-1 id of the uploaded content.

Possible HTTP status codes:

<table>
  <tr>
    <th>Status Code</th><th>Description</th>
  </tr>
  <tr>
    <td>201</td><td>Previously non-existent content was written successfully</td>
  </tr>
  <tr>
    <td>400</td><td>Client content is less than 1 byte or greater than 64MiB</td>
  </tr>
  <tr>
    <td>409</td><td>Content already exists on disk</td>
  </tr>
  <tr>
    <td>411</td><td>Content-length is required</td>
  </tr>
  <tr>
    <td>500</td><td>Generic server error</td>
  </tr>
  <tr>
    <td>507</td><td>Disk full</td>
  </tr>
</table>
