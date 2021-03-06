<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: event.proto

namespace Eventhouse\Grpc\Event\EntityStreamRequest;

use Google\Protobuf\Internal\GPBType;
use Google\Protobuf\Internal\RepeatedField;
use Google\Protobuf\Internal\GPBUtil;

/**
 * Generated from protobuf message <code>eventhouse.grpc.event.EntityStreamRequest.Filter</code>
 */
class Filter extends \Google\Protobuf\Internal\Message
{
    /**
     * Generated from protobuf field <code>string eventIdFrom = 1;</code>
     */
    protected $eventIdFrom = '';

    /**
     * Constructor.
     *
     * @param array $data {
     *     Optional. Data for populating the Message object.
     *
     *     @type string $eventIdFrom
     * }
     */
    public function __construct($data = NULL) {
        \GPBMetadata\Event::initOnce();
        parent::__construct($data);
    }

    /**
     * Generated from protobuf field <code>string eventIdFrom = 1;</code>
     * @return string
     */
    public function getEventIdFrom()
    {
        return $this->eventIdFrom;
    }

    /**
     * Generated from protobuf field <code>string eventIdFrom = 1;</code>
     * @param string $var
     * @return $this
     */
    public function setEventIdFrom($var)
    {
        GPBUtil::checkString($var, True);
        $this->eventIdFrom = $var;

        return $this;
    }

}

// Adding a class alias for backwards compatibility with the previous class name.
class_alias(Filter::class, \Eventhouse\Grpc\Event\EntityStreamRequest_Filter::class);

