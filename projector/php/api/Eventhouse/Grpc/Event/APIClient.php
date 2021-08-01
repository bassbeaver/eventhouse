<?php
// GENERATED CODE -- DO NOT EDIT!

namespace Eventhouse\Grpc\Event;

/**
 */
class APIClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \Eventhouse\Grpc\Event\PushRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function Push(\Eventhouse\Grpc\Event\PushRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/eventhouse.grpc.event.API/Push',
        $argument,
        ['\Eventhouse\Grpc\Event\Event', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Eventhouse\Grpc\Event\GetRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function Get(\Eventhouse\Grpc\Event\GetRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/eventhouse.grpc.event.API/Get',
        $argument,
        ['\Eventhouse\Grpc\Event\Event', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Eventhouse\Grpc\Event\EntityStreamRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\ServerStreamingCall
     */
    public function EntityStream(\Eventhouse\Grpc\Event\EntityStreamRequest $argument,
      $metadata = [], $options = []) {
        return $this->_serverStreamRequest('/eventhouse.grpc.event.API/EntityStream',
        $argument,
        ['\Eventhouse\Grpc\Event\Event', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Eventhouse\Grpc\Event\GlobalStreamRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\ServerStreamingCall
     */
    public function GlobalStream(\Eventhouse\Grpc\Event\GlobalStreamRequest $argument,
      $metadata = [], $options = []) {
        return $this->_serverStreamRequest('/eventhouse.grpc.event.API/GlobalStream',
        $argument,
        ['\Eventhouse\Grpc\Event\Event', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Eventhouse\Grpc\Event\SubscribeGlobalStreamRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\ServerStreamingCall
     */
    public function SubscribeGlobalStream(\Eventhouse\Grpc\Event\SubscribeGlobalStreamRequest $argument,
      $metadata = [], $options = []) {
        return $this->_serverStreamRequest('/eventhouse.grpc.event.API/SubscribeGlobalStream',
        $argument,
        ['\Eventhouse\Grpc\Event\Event', 'decode'],
        $metadata, $options);
    }

}
