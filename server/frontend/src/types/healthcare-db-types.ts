import { z } from 'zod';

// Event types and validation
export const EventTypeEnum = z.enum([
  'fall',
  'abnormal_temperature',
  'temperature_reading',
  'motion_detected'
]);

export type EventType = z.infer<typeof EventTypeEnum>;

export const EventSchema = z.object({
  ID: z.number(),
  CreatedAt: z.string().datetime().transform((val: string | number | Date) => new Date(val)),
  UpdatedAt: z.string().datetime().transform((val: string | number | Date) => new Date(val)),
  DeletedAt: z.string().datetime().transform((val: string | number | Date) => new Date(val)).nullable(),
  device_id: z.number(),
  type: EventTypeEnum
});

export type Event = z.infer<typeof EventSchema>;

// Alert types and validation
export const AlertSchema = z.object({
  ID: z.number(),
  CreatedAt: z.string().datetime().transform((val: string | number | Date) => new Date(val)),
  UpdatedAt: z.string().datetime().transform((val: string | number | Date) => new Date(val)),
  DeletedAt: z.string().datetime().transform((val: string | number | Date) => new Date(val)).nullable(),
  event_id: z.number().nullable(),
  event: EventSchema.nullable(),
  patient_id: z.number().nullable(),
  severity: z.string(),
  message: z.string(),
  resolved: z.boolean(),
  resolved_at: z.string().datetime().transform((val: string | number | Date) => new Date(val)).nullable()
});

export type Alert = z.infer<typeof AlertSchema>;

// Array schemas for list responses
export const AlertArraySchema = z.array(AlertSchema);
export type AlertArray = z.infer<typeof AlertArraySchema>;

export const EventArraySchema = z.array(EventSchema);
export type EventArray = z.infer<typeof EventArraySchema>;

// Parse single items
export function parseAlert(data: unknown): Alert {
  return AlertSchema.parse(data);
}

export function parseAlerts(data: unknown): AlertArray {
  return AlertArraySchema.parse(data);
}

export function parseEvent(data: unknown): Event {
  return EventSchema.parse(data);
}

export function parseEvents(data: unknown): EventArray {
  return EventArraySchema.parse(data);
}

// Safe parsing with error handling
export type ParseResult<T> = 
  | { success: true; data: T }
  | { success: false; error: string };

export function tryParseAlert(data: unknown): ParseResult<Alert> {
  try {
    return { success: true, data: parseAlert(data) };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error'
    };
  }
}

export function tryParseAlerts(data: unknown): ParseResult<AlertArray> {
  try {
    return { success: true, data: parseAlerts(data) };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error'
    };
  }
}

export function tryParseEvent(data: unknown): ParseResult<Event> {
  try {
    return { success: true, data: parseEvent(data) };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error'
    };
  }
}

export function tryParseEvents(data: unknown): ParseResult<EventArray> {
  try {
    return { success: true, data: parseEvents(data) };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error'
    };
  }
}

export const wsMessageSchema = z.preprocess(
  (input) => {
    if (typeof input === 'string') {
      try {
        return JSON.parse(input);
      } catch {
        return input;
      }
    }
    return input;
  },
  z.object({
    topic: z.string(),
    payload: z.unknown()
  })
)

export type WSMessage = z.infer<typeof wsMessageSchema>;

export function tryParseWSMessage(data: unknown): ParseResult<WSMessage> {
  try {
    return { success: true, data: wsMessageSchema.parse(data) };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error'
    };
  }
}