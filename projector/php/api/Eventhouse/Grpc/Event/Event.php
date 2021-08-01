<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: event.proto

namespace Eventhouse\Grpc\Event;

use Google\Protobuf\Internal\GPBType;
use Google\Protobuf\Internal\RepeatedField;
use Google\Protobuf\Internal\GPBUtil;

/**
 * Generated from protobuf message <code>eventhouse.grpc.event.Event</code>
 */
class Event extends \Google\Protobuf\Internal\Message
{
    /**
     * Generated from protobuf field <code>string eventId = 1;</code>
     */
    protected $eventId = '';
    /**
     * Generated from protobuf field <code>string eventType = 2;</code>
     */
    protected $eventType = '';
    /**
     * Generated from protobuf field <code>string entityType = 3;</code>
     */
    protected $entityType = '';
    /**
     * Generated from protobuf field <code>string entityId = 4;</code>
     */
    protected $entityId = '';
    /**
     * Generated from protobuf field <code>.google.protobuf.Timestamp recorded = 5;</code>
     */
    protected $recorded = null;
    /**
     * Generated from protobuf field <code>string payload = 6;</code>
     */
    protected $payload = '';

    /**
     * Constructor.
     *
     * @param array $data {
     *     Optional. Data for populating the Message object.
     *
     *     @type string $eventId
     *     @type string $eventType
     *     @type string $entityType
     *     @type string $entityId
     *     @type \Google\Protobuf\Timestamp $recorded
     *     @type string $payload
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

    /**
     * Generated from protobuf field <code>string eventType = 2;</code>
     * @return string
     */
    public function getEventType()
    {
        return $this->eventType;
    }

    /**
     * Generated from protobuf field <code>string eventType = 2;</code>
     * @param string $var
     * @return $this
     */
    public function setEventType($var)
    {
        GPBUtil::checkString($var, True);
        $this->eventType = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>string entityType = 3;</code>
     * @return string
     */
    public function getEntityType()
    {
        return $this->entityType;
    }

    /**
     * Generated from protobuf field <code>string entityType = 3;</code>
     * @param string $var
     * @return $this
     */
    public function setEntityType($var)
    {
        GPBUtil::checkString($var, True);
        $this->entityType = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>string entityId = 4;</code>
     * @return string
     */
    public function getEntityId()
    {
        return $this->entityId;
    }

    /**
     * Generated from protobuf field <code>string entityId = 4;</code>
     * @param string $var
     * @return $this
     */
    public function setEntityId($var)
    {
        GPBUtil::checkString($var, True);
        $this->entityId = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>.google.protobuf.Timestamp recorded = 5;</code>
     * @return \Google\Protobuf\Timestamp|null
     */
    public function getRecorded()
    {
        return isset($this->recorded) ? $this->recorded : null;
    }

    public function hasRecorded()
    {
        return isset($this->recorded);
    }

    public function clearRecorded()
    {
        unset($this->recorded);
    }

    /**
     * Generated from protobuf field <code>.google.protobuf.Timestamp recorded = 5;</code>
     * @param \Google\Protobuf\Timestamp $var
     * @return $this
     */
    public function setRecorded($var)
    {
        GPBUtil::checkMessage($var, \Google\Protobuf\Timestamp::class);
        $this->recorded = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>string payload = 6;</code>
     * @return string
     */
    public function getPayload()
    {
        return $this->payload;
    }

    /**
     * Generated from protobuf field <code>string payload = 6;</code>
     * @param string $var
     * @return $this
     */
    public function setPayload($var)
    {
        GPBUtil::checkString($var, True);
        $this->payload = $var;

        return $this;
    }

}
