<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: event.proto

namespace Eventhouse\Grpc\Event;

use Google\Protobuf\Internal\GPBType;
use Google\Protobuf\Internal\RepeatedField;
use Google\Protobuf\Internal\GPBUtil;

/**
 * Generated from protobuf message <code>eventhouse.grpc.event.GetRequest</code>
 */
class GetRequest extends \Google\Protobuf\Internal\Message
{
    /**
     * Generated from protobuf field <code>string eventId = 1;</code>
     */
    protected $eventId = '';

    /**
     * Constructor.
     *
     * @param array $data {
     *     Optional. Data for populating the Message object.
     *
     *     @type string $eventId
     * }
     */
    public function __construct($data = NULL) {
        \GPBMetadata\Event::initOnce();
        parent::__construct($data);
    }

    /**
     * Generated from protobuf field <code>string eventId = 1;</code>
     * @return string
     */
    public function getEventId()
    {
        return $this->eventId;
    }

    /**
     * Generated from protobuf field <code>string eventId = 1;</code>
     * @param string $var
     * @return $this
     */
    public function setEventId($var)
    {
        GPBUtil::checkString($var, True);
        $this->eventId = $var;

        return $this;
    }

}

