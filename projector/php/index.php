<?php

require __DIR__ . '/vendor/autoload.php';

const GRPC_API_URL = 'app:750';
const AUTH_HEADER = 'Y2xpZW50MTpzZWNyZXQx';
const EVENT_UID = '16278462494555805905';

$client = new \Eventhouse\Grpc\Event\APIClient(GRPC_API_URL, ['credentials' => Grpc\ChannelCredentials::createInsecure()]);

$getRequest = new \Eventhouse\Grpc\Event\GetRequest();
$getRequest->setEventId(EVENT_UID);

list($response, $status) = $client->Get($getRequest, ['Authorization' => ['Basic ' . AUTH_HEADER]])->wait();

if ($response instanceof \Eventhouse\Grpc\Event\Event) {
    echo \sprintf(
        " entity type: %s \n entity id: %s \n event type: %s \n event id: %s \n recorded: %s \n payload: %s",
        $response->getEntityType(),
        $response->getEntityId(),
        $response->getEventType(),
        $response->getEventId(),
        $response->getRecorded() ? $response->getRecorded()->toDateTime()->format('Y-m-d H:i:s.u') : 'null',
        $response->getPayload()
    );
} else {
    echo '$response is not type of ' . \Eventhouse\Grpc\Event\Event::class;
}

echo PHP_EOL . '$status: ' . \print_r($status, true);