(ns cas.storage
  (:gen-class))

(def ^:private store (ref {}))

(defn exists-kv? [key] (contains? @store key))

(defn read-kv [key] (get @store key nil))

(defn write-kv [key val] (dosync (ref-set store (assoc @store key val))))
